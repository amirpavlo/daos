#define DD_SUBSYS	DD_FAC(tests)

#include <getopt.h>
#include "daos_test.h"


char USAGE[] = "ds_cross_conn_test <warm_tier_group> <cold_tier_group>";
char TIER_1_ID[] = "TIER_1_ID";
char TIER_2_ID[] = "TIER_2_ID";

/*Small helpder func*/
static int
pool_create(const char *grp_id, uuid_t pool_id, daos_rank_list_t *svc)
{
	int rc = 0;

	rc = daos_pool_create(0731, geteuid(), getegid(), grp_id,
				      NULL, "pmem", 256 << 22, svc,
				      pool_id, NULL);
	return rc;
}

int
main(int argc, char **argv)
{
	int rank;
	int size;
	int rc;

	int warm_id_len;
	int cold_id_len;

	char *warm_grp;
	char *cold_grp;

	daos_handle_t	eqh;
	daos_event_t	ev;
	daos_event_t	*evp;



	/*IDs and daos (not MPI) rank info*/
	uuid_t warm_uuid;
	daos_rank_t warm_ranks[1];
	daos_rank_list_t warm_svc;

	uuid_t cold_uuid;
	daos_rank_t cold_ranks[1];
	daos_rank_list_t cold_svc;

	warm_svc.rl_nr.num = 1;
	warm_svc.rl_nr.num_out = 0;
	warm_svc.rl_ranks = warm_ranks;

	cold_svc.rl_nr.num = 1;
	cold_svc.rl_nr.num_out = 0;
	cold_svc.rl_ranks = cold_ranks;

	/*Pool handles and info*/
	daos_handle_t warm_poh;
	daos_pool_info_t warm_pool_info;

	if (argc != 3) {
		print_message("Incorrect number of args. %s\n", USAGE);
		return -1;
	}

	warm_id_len = strlen(argv[1]);
	cold_id_len = strlen(argv[2]);

	warm_grp = malloc(warm_id_len + 1);
	cold_grp = malloc(cold_id_len + 1);

	strcpy(warm_grp, argv[1]);
	strcpy(cold_grp, argv[2]);

	print_message("Warm-Tier Group: %s\n", warm_grp);
	print_message("Cold-Tier Group: %s\n", cold_grp);

	MPI_Init(&argc, &argv);

	MPI_Comm_rank(MPI_COMM_WORLD, &rank);
	MPI_Comm_size(MPI_COMM_WORLD, &size);
	MPI_Barrier(MPI_COMM_WORLD);

	rc = daos_init();

	if (rc) {
		print_message("daos_init() failed with %d\n", rc);
		return -1;
	}

	/*Initialze event and queue*/
	rc = daos_eq_create(&eqh);
	if (rc) {
		print_message("EQ Create Failed");
		D_GOTO(out, rc);
	}

	rc = daos_event_init(&ev, eqh, NULL);

	if (rank == 0) {
		/* create pools with minimal size */
		rc = pool_create(warm_grp, warm_uuid, &warm_svc);

		if (rc) {
			print_message("Warm Pool Create Failed: %d\n", rc);
			goto out;
		} else {
			print_message("Warm Pool Created\n");
		}
		/* create pools with minimal size */
		rc = pool_create(cold_grp, cold_uuid, &cold_svc);

		if (rc) {
			print_message("Cold Pool Create Failed: %d\n", rc);
			goto out;
		} else {
			print_message("Cold Pool Created\n");
		}

	}

	print_message("Warm Pool UUID: "DF_UUIDF"\n", DP_UUID(warm_uuid));
	print_message("Cold Pool UUID: "DF_UUIDF"\n", DP_UUID(cold_uuid));

	/*Register the colder tier ID and groupo server side*/
	print_message("Registering Cold Tier...\n");
	daos_tier_register_cold(cold_uuid, cold_grp, warm_uuid, warm_grp, NULL);

	print_message("Initiating Cross-Connect\n");
	/*Connect to the pool*/

	rc =  daos_tier_pool_connect(warm_uuid, warm_grp, &warm_svc,
				     DAOS_PC_RW, &warm_poh, &warm_pool_info,
				     &ev);

	print_message("Polling daos tier pool connect ev\n");
	daos_eq_poll(eqh, 1, DAOS_EQ_WAIT, 1, &evp);

	rc = evp->ev_error;

	if (rc)
		print_message("Pool Connect Failed with code: %d\n", rc);
	else
		print_message("Connected to pool\n");


	print_message("Disconnecting from Pools\n");
	rc = daos_pool_disconnect(warm_poh, NULL);
	if (rc)
		print_message("Warm Pool (local) Disconnect Failed\n");
	else
		print_message("Warm Pool (Local) Disconnect Success\n");

	print_message("Destroying Pools...\n");
	if (rank == 0) {
		rc = daos_pool_destroy(warm_uuid,
				       warm_grp /*grp*/, 1, NULL);

		if (rc)
			print_message("Error on warm pool destroy: %d\n", rc);
		else
			print_message("warm pool Destroyed\n");
		rc = daos_pool_destroy(cold_uuid,
				       cold_grp, 1, NULL);

		if (rc)
			print_message("Error on cold pool destroy: %d\n", rc);
		else
			print_message("Cold pool Destroyed\n");
	}

out:
	rc = daos_fini();
	if (rc)
		print_message("daos_fini() failed with %d\n", rc);

	MPI_Finalize();

	free(warm_grp);
	free(cold_grp);

	return 0;
}