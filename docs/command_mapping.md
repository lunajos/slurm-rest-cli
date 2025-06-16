# Slurm Traditional CLI to `srest` Command Mapping

This document maps traditional Slurm command-line tools to their equivalent `srest` commands. The `srest` tool provides a unified interface to the Slurm REST API while maintaining familiar command syntax.

## Command Structure

The `srest` command follows this general structure:

```
srest [global-options] <command> [subcommand] [options]
```

Where:
- `global-options`: Options that apply to all commands (e.g., `--token`, `--url`, `--format`)
- `command`: Primary command (e.g., `job`, `node`, `partition`)
- `subcommand`: Action to perform (e.g., `submit`, `show`, `update`)
- `options`: Command-specific options

## Command Mapping

### Job Management

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sbatch [options] script.sh` | `srest job submit [options] script.sh` | Submit a batch job |
| `salloc [options]` | `srest job allocate [options]` | Allocate resources for interactive job |
| `squeue [options]` | `srest job list [options]` | List jobs in queue |
| `scontrol show job <job_id>` | `srest job show <job_id>` | Show job details |
| `scontrol update job=<job_id> ...` | `srest job update <job_id> [options]` | Update job parameters |
| `scancel [options] <job_id>` | `srest job cancel [options] <job_id>` | Cancel a job |
| `scancel --signal=<signal> <job_id>` | `srest job signal <job_id> --signal=<signal>` | Send signal to a job |
| `sacct -j <job_id> [options]` | `srest job history <job_id> [options]` | Show job accounting information |

### Node Management

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sinfo [options]` | `srest node list [options]` | List node information |
| `scontrol show node <node_name>` | `srest node show <node_name>` | Show node details |
| `scontrol update nodename=<node_name> ...` | `srest node update <node_name> [options]` | Update node properties |
| `scontrol delete nodename=<node_name>` | `srest node delete <node_name>` | Delete a node |

### Partition Management

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sinfo -p [options]` | `srest partition list [options]` | List partition information |
| `scontrol show partition <partition_name>` | `srest partition show <partition_name>` | Show partition details |

### Reservation Management

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `scontrol show reservation` | `srest reservation list` | List all reservations |
| `scontrol show reservation <name>` | `srest reservation show <name>` | Show reservation details |

### QOS Management

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sacctmgr show qos [options]` | `srest qos list [options]` | List QOS information |
| `sacctmgr show qos <name>` | `srest qos show <name>` | Show QOS details |
| `sacctmgr delete qos <name>` | `srest qos delete <name>` | Delete QOS |

### TRES Management

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sacctmgr show tres` | `srest tres list` | List TRES information |
| `sacctmgr add tres ...` | `srest tres add [options]` | Add TRES |

### Configuration

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sacctmgr show configuration` | `srest config show` | Show configuration |
| `sacctmgr modify configuration ...` | `srest config update [options]` | Update configuration |

### Diagnostics and System Operations

| Traditional Command | New `srest` Command | Description |
|--------------------|---------------------|-------------|
| `sdiag` | `srest diag` | Show diagnostic information |
| `scontrol ping` | `srest ping` | Ping Slurm controller |
| `scontrol show licenses` | `srest license list` | Show license information |
| `scontrol reconfigure` | `srest reconfigure` | Reconfigure Slurm |

## Global Options

All `srest` commands support these global options:

| Option | Description |
|--------|-------------|
| `--token=<token>` | Authentication token |
| `--jwt=<jwt>` | JWT for authentication |
| `--url=<url>` | Slurm REST API URL |
| `--user=<username>` | Username for authentication |
| `--format=<format>` | Output format (json, yaml, table) |
| `--curl` | Show equivalent curl command instead of executing |
| `--help` | Show help for command |
| `--version` | Show version information |

## Examples

### Job Submission

Traditional:
```
sbatch --job-name=test --nodes=2 --ntasks=4 --time=1:00:00 script.sh
```

New:
```
srest job submit --job-name=test --nodes=2 --ntasks=4 --time=1:00:00 script.sh
```

### Viewing Jobs

Traditional:
```
squeue --user=$USER
```

New:
```
srest job list --user=$USER
```

### Canceling Jobs

Traditional:
```
scancel 12345
```

New:
```
srest job cancel 12345
```

### Node Information

Traditional:
```
sinfo -N
```

New:
```
srest node list
```

## Notes

1. The `srest` command maintains backward compatibility with traditional Slurm options.
2. All commands support the `--help` flag for detailed usage information.
3. The `--curl` option allows users to see the equivalent curl command for API calls.
4. Minimum unique prefix for subcommands is supported (e.g., `srest j s` for `srest job submit`).
