#!/usr/bin/python
'''
  (C) Copyright 2018-2019 Intel Corporation.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
  The Government's rights to use, modify, reproduce, release, perform, display,
  or disclose this software are subject to the terms of the Apache License as
  provided in Contract No. B609815.
  Any reproduction of computer software, computer software documentation, or
  portions thereof marked with this legend must also reproduce the markings.
'''
from __future__ import print_function

import os
import json
import distutils.spawn
from apricot import Test

import agent_utils
import server_utils
import write_host_file
import ior_utils
from daos_api import DaosContext, DaosPool, DaosApiError

class EightServers(Test):
    """
    Test class Description: Runs IOR with 8 servers.
    """

    def __init__(self, *args, **kwargs):

        super(EightServers, self).__init__(*args, **kwargs)

        self.basepath = None
        self.server_group = None
        self.context = None
        self.pool = None
        self.slots = None
        self.hostlist_servers = None
        self.hostfile_servers = None
        self.hostlist_clients = None
        self.hostfile_clients = None

    def setUp(self):
        self.agent_sessions = None
        # get paths from the build_vars generated by build
        with open('../../../.build_vars.json') as build_file:
            build_paths = json.load(build_file)
        self.basepath = os.path.normpath(build_paths['PREFIX'] + "/../")

        self.server_group = self.params.get("name", '/server_config/',
                                            'daos_server')

        # setup the DAOS python API
        self.context = DaosContext(build_paths['PREFIX'] + '/lib/')

        self.hostlist_servers = self.params.get("test_servers",
                                                '/run/hosts/test_machines/*')
        self.hostfile_servers = (
            write_host_file.write_host_file(self.hostlist_servers,
                                            self.workdir))
        print("Host file servers is: {}".format(self.hostfile_servers))

        self.hostlist_clients = self.params.get("test_clients",
                                                '/run/hosts/test_machines/*')
        self.slots = self.params.get("slots", '/run/ior/clientslots/*')
        self.hostfile_clients = (
            write_host_file.write_host_file(self.hostlist_clients, self.workdir,
                                            self.slots))
        print("Host file clients is: {}".format(self.hostfile_clients))

        self.agent_sessions = agent_utils.run_agent(self.basepath,
                                                    self.hostlist_servers,
                                                    self.hostlist_clients)
        server_utils.run_server(self.hostfile_servers, self.server_group,
                                self.basepath)

        if not distutils.spawn.find_executable("ior") and \
           int(str(self.name).split("-")[0]) == 1:
            ior_utils.build_ior(self.basepath)

    def tearDown(self):
        try:
            if self.pool is not None and self.pool.attached:
                self.pool.destroy(1)
        finally:
            if self.agent_sessions:
                agent_utils.stop_agent(self.agent_sessions,
                                       self.hostlist_clients)
            server_utils.stop_server(hosts=self.hostlist_servers)

    def executable(self, iorflags=None):
        """
        Executable function to run ior for sequential and random order
        """

        # parameters used in pool create
        createmode = self.params.get("mode", '/run/pool/createmode/*/')
        createuid = os.geteuid()
        creategid = os.getegid()
        createsetid = self.params.get("setname", '/run/pool/createset/')
        createsize = self.params.get("size", '/run/pool/createsize/')
        createsvc = self.params.get("svcn", '/run/pool/createsvc/')
        iteration = self.params.get("iter", '/run/ior/iteration/')
        block_size = self.params.get("blocksize", '/run/ior/clientslots/*')
        record_size = self.params.get("r", '/run/ior/recordsize/*')
        stripe_count = self.params.get("c", '/run/ior/stripecount/')
        async_io = self.params.get("a", '/run/ior/asyncio/')
        object_class = self.params.get("o", '/run/ior/objectclass/*/')
        transfer_size = self.params.get("t",
                                        '/run/ior/transfersize_stripesize/*/')
        stripe_size = self.params.get("s",
                                      '/run/ior/transfersize_stripesize/*/')
        expected_result = 'PASS'

        if record_size == '4k' and transfer_size == '1k':
            expected_result = 'FAIL'

        try:
            # initialize a python pool object then create the underlying
            # daos storage
            self.pool = DaosPool(self.context)
            self.pool.create(createmode, createuid, creategid,
                             createsize, createsetid, None, None, createsvc)

            pool_uuid = self.pool.get_uuid_str()
            svc_list = ""
            for item in range(createsvc):
                svc_list += str(int(self.pool.svc.rl_ranks[item])) + ":"
            svc_list = svc_list[:-1]

            print ("svc_list: {}".format(svc_list))

            ior_utils.run_ior(self.hostfile_clients, iorflags, iteration,
                              block_size, transfer_size, pool_uuid, svc_list,
                              record_size, stripe_size, stripe_count,
                              async_io, object_class, self.basepath, self.slots)

            if expected_result == 'FAIL':
                self.fail("Test was expected to fail but it passed.\n")

        except (DaosApiError, ior_utils.IorFailed) as excep:
            print(excep)
            if expected_result != 'FAIL':
                self.fail("Test was expected to pass but it failed.\n")

    def test_sequential(self):
        """
        Test ID: DAOS-1264
        Test Description: Run IOR with 32,64 and 128 clients config
                          sequentially.
        Use Cases: Different combinations of 32/64/128 Clients, 8b/1k/4k
                   record size, 1k/4k/1m/8m transfersize and stripesize
                   and 16 async io.
        :avocado: tags=ior,eightservers,ior_sequential
        """
        ior_flags = self.params.get("F", '/run/ior/iorflags/sequential/')
        self.executable(ior_flags)

    def test_random(self):
        """
        Test ID: DAOS-1264
        Test Description: Run IOR with 32,64 and 128 clients config in random
                          order.
        Use Cases: Different combinations of 32/64/128 Clients, 8b/1k/4k
                   record size, 1k/4k/1m/8m transfersize and stripesize
                   and 16 async io.
        :avocado: tags=ior,eightservers,ior_random
        """
        ior_flags = self.params.get("F", '/run/ior/iorflags/random/')
        self.executable(ior_flags)
