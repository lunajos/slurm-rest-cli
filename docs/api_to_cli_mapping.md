# Slurm REST API to CLI Command Mapping

This document provides a comprehensive mapping between the Slurm REST API (v0.0.42) endpoints/parameters and their corresponding command line interface (CLI) commands and flags.

## Authentication Methods

The Slurm REST API supports three authentication methods:

1. **User + Token headers**: `X-SLURM-USER-NAME` + `X-SLURM-USER-TOKEN`
2. **Token only**: `X-SLURM-USER-TOKEN`
3. **Bearer Authentication**: JWT

CLI equivalent authentication is typically handled via:
- Environment variables: `SLURM_TOKEN`, `SLURM_JWT`
- Command line options: `--token`, `--jwt`

## Job Submission and Management

### Job Submission (`/slurm/v0.0.42/job/submit`)

| API Parameter | CLI Command | CLI Flag/Option |
|---------------|-------------|----------------|
| `script` (deprecated) | `sbatch` | Script file as argument |
| `job.name` | `sbatch` | `--job-name=<name>` or `-J <name>` |
| `job.account` | `sbatch` | `--account=<account>` or `-A <account>` |
| `job.partition` | `sbatch` | `--partition=<partition>` or `-p <partition>` |
| `job.qos` | `sbatch` | `--qos=<qos>` |
| `job.comment` | `sbatch` | `--comment=<comment>` |
| `job.nodes` | `sbatch` | `--nodes=<count>` or `-N <count>` |
| `job.tasks` | `sbatch` | `--ntasks=<count>` or `-n <count>` |
| `job.cpus_per_task` | `sbatch` | `--cpus-per-task=<count>` or `-c <count>` |
| `job.mem_per_cpu` | `sbatch` | `--mem-per-cpu=<MB>` |
| `job.mem_per_node` | `sbatch` | `--mem=<MB>` |
| `job.time_limit` | `sbatch` | `--time=<time>` or `-t <time>` |
| `job.begin_time` | `sbatch` | `--begin=<time>` |
| `job.deadline` | `sbatch` | `--deadline=<time>` |
| `job.std_in` | `sbatch` | `--input=<file>` or `-i <file>` |
| `job.std_out` | `sbatch` | `--output=<file>` or `-o <file>` |
| `job.std_err` | `sbatch` | `--error=<file>` or `-e <file>` |
| `job.hold` | `sbatch` | `--hold` or `-H` |
| `job.requeue` | `sbatch` | `--requeue` or `--no-requeue` |
| `job.kill_on_node_fail` | `sbatch` | `--kill-on-invalid-dep=<yes\|no>` |
| `job.array_inx` | `sbatch` | `--array=<indexes>` or `-a <indexes>` |
| `job.dependency` | `sbatch` | `--dependency=<dependency_list>` or `-d <dependency_list>` |
| `job.mail_type` | `sbatch` | `--mail-type=<type>` |
| `job.mail_user` | `sbatch` | `--mail-user=<user>` |
| `job.nice` | `sbatch` | `--nice=<adjustment>` |
| `job.constraints` | `sbatch` | `--constraint=<list>` or `-C <list>` |
| `job.x11` | `sbatch` | `--x11=<mode>` |
| `job.gres` | `sbatch` | `--gres=<list>` |
| `job.tres_per_job` | `sbatch` | `--tres-per-job=<list>` |
| `job.tres_per_node` | `sbatch` | `--tres-per-node=<list>` |
| `job.tres_per_socket` | `sbatch` | `--tres-per-socket=<list>` |
| `job.tres_per_task` | `sbatch` | `--tres-per-task=<list>` |
| `job.cpus_per_tres` | `sbatch` | `--cpus-per-gpu=<count>` |
| `job.mem_per_tres` | `sbatch` | `--mem-per-gpu=<MB>` |
| `job.licenses` | `sbatch` | `--licenses=<license>` or `-L <license>` |
| `job.clusters` | `sbatch` | `--clusters=<list>` or `-M <list>` |
| `job.reservation` | `sbatch` | `--reservation=<name>` |
| `job.priority` | `sbatch` | `--priority=<value>` |
| `job.current_working_directory` | `sbatch` | `--chdir=<directory>` or `-D <directory>` |
| `job.workdir` | `sbatch` | `--workdir=<directory>` |
| `job.wckey` | `sbatch` | `--wckey=<wckey>` |
| `job.exclusive` | `sbatch` | `--exclusive` or `--exclusive=user` |
| `job.shared` | `sbatch` | `--share` or `--no-share` |
| `job.oversubscribe` | `sbatch` | `--oversubscribe` or `--exclusive` |
| `job.contiguous` | `sbatch` | `--contiguous` |
| `job.core_spec` | `sbatch` | `--core-spec=<count>` |
| `job.thread_spec` | `sbatch` | `--thread-spec=<count>` |
| `job.min_cpus` | `sbatch` | `--mincpus=<count>` |
| `job.min_nodes` | `sbatch` | `--nodes=<min[-max]>` or `-N <min[-max]>` |
| `job.max_nodes` | `sbatch` | `--nodes=<min-max>` or `-N <min-max>` |
| `job.sockets_per_node` | `sbatch` | `--sockets-per-node=<count>` |
| `job.cores_per_socket` | `sbatch` | `--cores-per-socket=<count>` |
| `job.threads_per_core` | `sbatch` | `--threads-per-core=<count>` |
| `job.ntasks_per_node` | `sbatch` | `--ntasks-per-node=<count>` |
| `job.ntasks_per_socket` | `sbatch` | `--ntasks-per-socket=<count>` |
| `job.ntasks_per_core` | `sbatch` | `--ntasks-per-core=<count>` |
| `job.ntasks_per_tres` | `sbatch` | `--ntasks-per-gpu=<count>` |
| `job.environment` | `sbatch` | `--export=<environment_variables>` |
| `job.burst_buffer` | `sbatch` | `--bb=<spec>` |
| `job.delay_boot` | `sbatch` | `--delay-boot=<minutes>` |
| `job.network` | `sbatch` | `--network=<type>` |

### Job Allocation (`/slurm/v0.0.42/job/allocate`)

| API Parameter | CLI Command | CLI Flag/Option |
|---------------|-------------|----------------|
| Similar to job submission | `salloc` | Same flags as `sbatch` |

### Job Query (`/slurm/v0.0.42/jobs/`)

| API Parameter | CLI Command | CLI Flag/Option |
|---------------|-------------|----------------|
| `update_time` | `squeue` | `--starttime=<time>` |
| `flags` (ALL, DETAIL, etc.) | `squeue` | Various flags like `--all`, `--details` |
| N/A | `squeue` | `-o <format>` or `--format=<format>` |

### Job Control (`/slurm/v0.0.42/job/{job_id}`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol show job <job_id>` | N/A |
| POST (update) | `scontrol update job=<job_id> ...` | Various parameters |
| DELETE (cancel) | `scancel <job_id>` | N/A |
| DELETE with signal | `scancel --signal=<signal> <job_id>` | `--signal=<signal>` |
| DELETE with flags | `scancel` | `--batch`, `--full`, etc. |

## Node Management

### Node Query (`/slurm/v0.0.42/nodes/`)

| API Parameter | CLI Command | CLI Flag/Option |
|---------------|-------------|----------------|
| `update_time` | `sinfo` | N/A |
| `flags` | `sinfo` | Various flags |
| N/A | `sinfo` | `-o <format>` or `--format=<format>` |

### Node Control (`/slurm/v0.0.42/node/{node_name}`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol show node <node_name>` | N/A |
| POST (update) | `scontrol update nodename=<node_name> ...` | Various parameters |
| DELETE | `scontrol delete nodename=<node_name>` | N/A |

## Partition Management

### Partition Query (`/slurm/v0.0.42/partitions/`)

| API Parameter | CLI Command | CLI Flag/Option |
|---------------|-------------|----------------|
| `update_time` | `sinfo` | N/A |
| `flags` | `sinfo` | Various flags |
| N/A | `sinfo -p` | `-p <partition>` or `--partition=<partition>` |

### Partition Control (`/slurm/v0.0.42/partition/{partition_name}`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol show partition <partition_name>` | N/A |

## Reservation Management

### Reservation Query (`/slurm/v0.0.42/reservations/`)

| API Parameter | CLI Command | CLI Flag/Option |
|---------------|-------------|----------------|
| `update_time` | `scontrol show reservation` | N/A |

### Reservation Control (`/slurm/v0.0.42/reservation/{reservation_name}`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol show reservation <reservation_name>` | N/A |

## SlurmDB Operations

### Job History (`/slurmdb/v0.0.42/job/{job_id}`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `sacct -j <job_id>` | `-j <job_id>` or `--jobs=<job_id>` |

### Configuration (`/slurmdb/v0.0.42/config`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `sacctmgr show configuration` | N/A |
| POST | `sacctmgr modify configuration ...` | Various parameters |

### TRES Management (`/slurmdb/v0.0.42/tres/`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `sacctmgr show tres` | N/A |
| POST | `sacctmgr add tres ...` | Various parameters |

### QOS Management (`/slurmdb/v0.0.42/qos/`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `sacctmgr show qos` | N/A |
| GET with filters | `sacctmgr show qos <name>` | Various filters |
| DELETE | `sacctmgr delete qos <name>` | N/A |

## Other Operations

### Diagnostics (`/slurm/v0.0.42/diag`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `sdiag` | N/A |

### Ping (`/slurm/v0.0.42/ping`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol ping` | N/A |

### Licenses (`/slurm/v0.0.42/licenses`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol show licenses` | N/A |

### Reconfigure (`/slurm/v0.0.42/reconfigure`)

| API Method | CLI Command | CLI Flag/Option |
|------------|-------------|----------------|
| GET | `scontrol reconfigure` | N/A |

## Notes

1. This mapping is based on Slurm REST API v0.0.42 and may need updates for other versions.
2. Some CLI commands have additional options not directly mapped to API parameters.
3. The REST API may have additional capabilities not available in the CLI and vice versa.
4. Authentication methods in the CLI may vary based on the implementation.
