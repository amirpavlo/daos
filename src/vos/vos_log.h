/**
 * (C) Copyright 2019 Intel Corporation.
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
 * VOS punch log
 * vos/vos_log.h
 *
 * Author: Jeff Olivier <jeffrey.v.olivier@intel.com>
 */

#ifndef _VOS_LOG_H
#define _VOS_LOG_H
#include <daos_srv/vos_types.h>
#include <daos_srv/dtx_srv.h>

/** Number of log entries to allocate at a time */
#define VOS_LOG_SZ 8

/** The offset will always be aligned to malloc alignment.  Use smallest bit for
 *  punch flag
 */
#define VOS_PUNCH_BIT (1ULL)

struct vos_log_entry_df {
	/** timestamp of entry */
	uint64_t	le_epoch;
	/** punch flag and either dtx offset or count */
	uint64_t	le_data;
};

struct vos_log_df {
	/** offset of next set of log entries */
	uint64_t	l_next;
	/** offset of previous set of log entries */
	uint64_t	l_prev;
	/** block of entries.   To save on-disk space, it is assumed that
	 *  epoch == 0 marks an unused entry
	 */
	struct vos_log_entry_df	l_entries[VOS_LOG_SZ]
};

struct vos_log_root_df {
	/** Offset of log.  If it's 0, no log is present */
	uint64_t		lh_log;
	struct vos_log_entry_df	lh_entry;
};

/** Initialize the root for the first time.  Caller is expected to be inside
 *  a transaction and it is assumed that the struct has already been added
 *  to the log.
 *
 *  \param	root[IN,OUT]	The log root to initialize
 *  \param	epc[IN]		The epoch of the tx
 *  \param	dtx_off[IN]	The offset of the dtx in the pool
 *  \param	punch[IN]	Operation is a punch
 */
static inline
void vos_log_root_init(struct vos_log_root_df *root, daos_epoch_t epc,
		       uint64_t dtx_off, bool punch)
{
	root->lh_log = 0;
	root->lh_entry.le_epoch = epc;
	root->lh_entry.le_data = dtx_off;
	if (punch)
		root->lh_entry.le_data |= VOS_PUNCH_BIT;
}

#endif
