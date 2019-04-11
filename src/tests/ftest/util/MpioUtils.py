#!/usr/bin/python
'''
    (C) Copyright 2019 Intel Corporation.

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
    The Government's rights to use, modify, reproduce, release, perform,
    display, or disclose this software are subject to the terms of the Apache
    License as provided in Contract No. B609815.
    Any reproduction of computer software, computer software documentation, or
    portions thereof marked with this legend must also reproduce the markings.
    '''
from __future__ import print_function
import os
import subprocess
import paramiko
import socket

class MpioFailed(Exception):
    """Raise if MPIO failed"""

class MpioUtils():
    """MpioUtils Class"""

    def __init__(self):

        self.mpichinstall = None

    def mpich_installed(self, hostlist):
        """Check if mpich is installed"""

        try:
            # checking mpich install
            checkmpich = subprocess.check_output(["ssh", hostlist[0],
                                                  "which mpichversion"])

            # Obtaning the location where mpich is installed by
            # executing "mpichversion"
            self.mpichinstall = checkmpich[:-17]
            print (self.mpichinstall)
            return True

        except subprocess.CalledProcessError as error:
            print ("Mpich not installed \n {}".format(error))
            return False

    def run_romio(self, basepath, hostlist, romio_test_repo):
        """
            Running ROMIO testsuite under mpich
            Function Arguments:
                basepath        --path where all daos and it's dependencies can
                                  be fetched from
                hostlist        --list of client hosts
                romio_test_repo --built romio test directory
        """

        # environment variables only to be set on client node
        env_variables = ["export PATH={}/bin:$PATH".format(self.mpichinstall),
                         "export LD_LIBRARY_PATH={}/lib:$LD_LIBRARY_PATH"\
                         .format(self.mpichinstall),
                         "export CRT_ATTACH_INFO_PATH={}/install/tmp/"\
                         .format(basepath),
                         "export MPI_LIB=''", "export DAOS_SINGLETON_CLI=1"]

        # setting attributes
        cd_cmd = 'cd ' + romio_test_repo
        test_cmd = './runtests -fname=daos:test1 -subset -daos'
        run_cmd = cd_cmd + ' && ' + test_cmd
        print ("Romio test run command: {}".format(run_cmd))

        try:
            # establish conection and run romio test
            # All the tests must pass with "No Errors"
            # if any test fails it should return "non-zero exit code"
            ssh = paramiko.SSHClient()
            ssh.load_system_host_keys()
            ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
            ssh.connect(hostlist[0])
            ssh_stdin, ssh_stdout, ssh_stderr = ssh.exec_command(env_variables[0] \
                                                + " && " + env_variables[1] + \
                                                " && " + env_variables[2] + \
                                                " && " + env_variables[3] + \
                                                " && " + env_variables[4] + \
                                                " && " + run_cmd)
            print (ssh_stdout.read())
            print (ssh_stderr.read())

        except (IOError, OSError, paramiko.SSHException, socket.error) as error:
            raise MpioFailed("<ROMIO Test FAILED> \nException occurred: {}"\
                             .format(str(error)))

    def run_llnl_mpi4py(self, basepath, hostfile, pool_uuid, test_repo,
                        test_name):
        """
            Running LLNL and MPI4PY testsuites
            Function Arguments:
                basepath          --path where all daos and it's dependencies
                                    can be fetched from
                hostfile          --client hostfile
                pool_uuid         --Pool UUID
                test_repo         --test repo location
                test_name         --name of test to be tested
        """

        # environment variables only to be set on client node
        env_variables = ["export PATH={}/bin:$PATH".format(self.mpichinstall),
                         "export LD_LIBRARY_PATH={}/lib:$LD_LIBRARY_PATH"
                         .format(self.mpichinstall),
                         "export CRT_ATTACH_INFO_PATH={}/install/tmp/"
                         .format(basepath),
                         "export MPI_LIB=''", "export DAOS_SINGLETON_CLI=1",
                         "export MPIO_USER_PATH=daos:",
                         "export DAOS_POOL={}".format(pool_uuid),
                         "export DAOS_SVCL={}".format(0)]

        # setting attributes
        cd_cmd = 'cd ' + test_repo
        # running 8 client processes
        if test_name == "llnl" and os.path.isfile(test_repo + "/testmpio_daos"):
            test_cmd = "mpirun -np 8 --hostfile {} ./testmpio_daos 1"\
                           .format(hostfile)
        elif test_name == "mpi4py" and os.path.isfile(test_repo +
                                                      "/test_io_daos.py"):
            test_cmd = "mpiexec -n 8 python test_io_daos.py"
        else:
            raise MpioFailed("Wrong test name or test repo location specified")

        run_cmd = cd_cmd + ";" + env_variables[0] + ";" + env_variables[1] \
                  + ";" + env_variables[2] + ";" + env_variables[3] + ";" \
                  + env_variables[4] + ";" + env_variables[5] + ";" \
                  + env_variables[6] + ";" + env_variables[7] + ";" + test_cmd

        print ("run command: {}".format(run_cmd))

        try:
            process = subprocess.Popen(run_cmd, stdout=subprocess.PIPE,
                                       stderr=subprocess.STDOUT, shell=True)
            while True:
                output = process.stdout.readline()
                if output == '' and process.poll() is not None:
                    break
                if output:
                    print (output.strip())
            if process.poll() != 0:
                raise MpioFailed("{} Run process".format(test_name)
                                 + " Failed with non zero exit"
                                 + " code:{}".format(process.poll()))
        except (ValueError, OSError) as error:
            raise MpioFailed("<Test FAILED> \nException occurred: {}"\
                             .format(str(error)))
