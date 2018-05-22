/**
 * (C) Copyright 2017 Intel Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
 * The Government's rights to use, modify, reproduce, release, perform, display,
 * or disclose this software are subject to the terms of the Apache License as
 * provided in Contract No. B609815.
 * Any reproduction of computer software, computer software documentation, or
 * portions thereof marked with this legend must also reproduce the markings.
 */
/**
 * rdb: KVSs
 *
 * This file implements an LRU cache of rdb_kvs objects, each of which maps a
 * KVS path to the matching VOS object. The cache provides better KVS path
 * lookup performance.
 */

#define D_LOGFAC	DD_FAC(rdb)

#include <daos_srv/rdb.h>

#include "rdb_internal.h"
#include "rdb_layout.h"

struct rdb_kvs_open_arg {
	struct rdb     *deo_db;
	rdb_oid_t	deo_parent;
	uint64_t	deo_index;
};

/* Open key in arg->deo_parent. */
static int
rdb_kvs_open_path_cb(daos_iov_t *key, void *varg)
{
	struct rdb_kvs_open_arg	       *arg = varg;
	rdb_oid_t			parent = arg->deo_parent;
	daos_iov_t			value;

	if (key->iov_len == 0) {
		D_ASSERTF(parent == RDB_LC_ATTRS, DF_X64"\n", parent);
		key = &rdb_lc_root;
	}
	daos_iov_set(&value, &arg->deo_parent, sizeof(arg->deo_parent));
	return rdb_lc_lookup(arg->deo_db->d_lc, arg->deo_index, parent, key,
			     &value);
}

/*
 * Open the KVS corresponding to path, which is not in the cache. Currently,
 * the result is just an object ID, since object handles are not exported.
 */
static int
rdb_kvs_open_path(struct rdb *db, uint64_t index, const rdb_path_t *path,
		  rdb_oid_t *object)
{
	rdb_path_t		p = *path;
	struct rdb_kvs	       *kvs = NULL;
	struct rdb_kvs_open_arg	arg;
	int			rc;

	/* See if we can find a cache hit for a prefix of the path. */
	while (rdb_path_pop(&p) == 0 && p.iov_len > 0) {
		rc = rdb_kvs_lookup(db, &p, index, false /* alloc */, &kvs);
		if (rc == 0)
			break;
		else if (rc != -DER_NONEXIST)
			return rc;
	};

	/* Walk through the keys after "p". */
	D_DEBUG(DB_ANY, "walking path %zu from kvs %p\n", p.iov_len, kvs);
	p.iov_buf += p.iov_len;
	p.iov_buf_len -= p.iov_len;
	p.iov_len = path->iov_len - p.iov_len;
	D_ASSERT(p.iov_len > 0);
	arg.deo_db = db;
	arg.deo_parent = kvs == NULL ? RDB_LC_ATTRS : kvs->de_object;
	arg.deo_index = index;
	rc = rdb_path_iterate(&p, rdb_kvs_open_path_cb, &arg);
	if (kvs != NULL)
		rdb_kvs_put(db, kvs);
	if (rc != 0)
		return rc;

	D_DEBUG(DB_ANY, "got kvs handle "DF_X64"\n", arg.deo_parent);
	*object = arg.deo_parent;
	return 0;
}

static inline struct rdb_kvs *
rdb_kvs_obj(struct daos_llink *entry)
{
	return container_of(entry, struct rdb_kvs, de_entry);
}

struct rdb_kvs_alloc_arg {
	struct rdb     *dea_db;
	uint64_t	dea_index;
	bool		dea_alloc;
};

static int
rdb_kvs_alloc_ref(void *key, unsigned int ksize, void *varg,
		  struct daos_llink **link)
{
	struct rdb_kvs_alloc_arg       *arg = varg;
	struct rdb_kvs		       *kvs;
	void			       *buf;
	int				rc;

	if (!arg->dea_alloc) {
		rc = -DER_NONEXIST;
		goto err;
	}

	D_ALLOC_PTR(kvs);
	if (kvs == NULL)
		D_GOTO(err, rc = -DER_NOMEM);

	/* kvs->de_path */
	D_ALLOC(buf, ksize);
	if (buf == NULL)
		D_GOTO(err_kvs, rc = -DER_NOMEM);
	memcpy(buf, key, ksize);
	daos_iov_set(&kvs->de_path, buf, ksize);

	/* kvs->de_object */
	rc = rdb_kvs_open_path(arg->dea_db, arg->dea_index, &kvs->de_path,
			       &kvs->de_object);
	if (rc != 0)
		D_GOTO(err_path, rc);

	D_DEBUG(DB_ANY, DF_DB": created %p len %u\n", DP_DB(arg->dea_db), kvs,
		ksize);
	*link = &kvs->de_entry;
	return 0;

err_path:
	D_FREE(kvs->de_path.iov_buf);
err_kvs:
	D_FREE_PTR(kvs);
err:
	return rc;
}

static void
rdb_kvs_free_ref(struct daos_llink *llink)
{
	struct rdb_kvs *kvs = rdb_kvs_obj(llink);

	D_DEBUG(DB_ANY, "freeing %p "DF_X64"\n", kvs, kvs->de_object);
	D_FREE(kvs->de_path.iov_buf);
	D_FREE_PTR(kvs);
}

static bool
rdb_kvs_cmp_keys(const void *key, unsigned int ksize, struct daos_llink *llink)
{
	struct rdb_kvs *kvs = rdb_kvs_obj(llink);

	if (ksize != kvs->de_path.iov_len)
		return false;
	if (memcmp(key, kvs->de_path.iov_buf, ksize) != 0)
		return false;
	return true;
}

static struct daos_llink_ops rdb_kvs_cache_ops = {
	.lop_alloc_ref	= rdb_kvs_alloc_ref,
	.lop_free_ref	= rdb_kvs_free_ref,
	.lop_cmp_keys	= rdb_kvs_cmp_keys
};

int
rdb_kvs_cache_create(struct daos_lru_cache **cache)
{
	return daos_lru_cache_create(5 /* bits */, D_HASH_FT_NOLOCK /* feats */,
				     &rdb_kvs_cache_ops, cache);
}

void
rdb_kvs_cache_destroy(struct daos_lru_cache *cache)
{
	daos_lru_cache_destroy(cache);
}

void
rdb_kvs_cache_evict(struct daos_lru_cache *cache)
{
	daos_lru_cache_evict(cache, NULL /* cond */, NULL /* args */);
}

int
rdb_kvs_lookup(struct rdb *db, const rdb_path_t *path, uint64_t index,
	       bool alloc, struct rdb_kvs **kvs)
{
	struct rdb_kvs_alloc_arg	arg;
	struct daos_llink	       *entry;
	int				rc;

	D_DEBUG(DB_TRACE, DF_DB": looking up "DF_IOV": alloc=%d\n", DP_DB(db),
		DP_IOV(path), alloc);

	arg.dea_db = db;
	arg.dea_index = index;
	arg.dea_alloc = alloc;
	rc = daos_lru_ref_hold(db->d_kvss, path->iov_buf, path->iov_len, &arg,
			       &entry);
	if (rc != 0)
		return rc;

	*kvs = rdb_kvs_obj(entry);
	return 0;
}

void
rdb_kvs_put(struct rdb *db, struct rdb_kvs *kvs)
{
	daos_lru_ref_release(db->d_kvss, &kvs->de_entry);
}

void
rdb_kvs_evict(struct rdb *db, struct rdb_kvs *kvs)
{
	daos_lru_ref_evict(&kvs->de_entry);
}