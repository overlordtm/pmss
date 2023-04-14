# Useful queries

Do not forget to end with LIMIT when exploring, it can be quite a lot of matches

## Create a view to see latest report for each machine

```sql
CREATE VIEW latest_reports AS

WITH latest_reports AS (
    SELECT
        machine_id,
        MAX(created_at) AS latest_created_at
    FROM
        report_runs
    GROUP BY
        machine_id
)

SELECT
    m.ID,
    m.hostname,
    r.ID AS report_id,
    r.created_at AS report_created_at
FROM
    machines m
JOIN
    latest_reports lr ON m.ID = lr.machine_id
JOIN
    report_runs r ON r.machine_id = lr.machine_id AND r.created_at = lr.latest_created_at
```

## Get all machines that have files with given prefix

```sql
SELECT DISTINCT m.hostname FROM scanned_files f JOIN machines m ON f.machine_id = m.id WHERE f.path LIKE '/var/lib/mysql%'
```

## Find "rare" files

```sql
SELECT md5, COUNT(*) as cnt from scanned_files WHERE path NOT LIKE '/etc/%' AND path NOT LIKE '/var/lib/docker/%' AND path NOT LIKE '/var/lib/kubelet/pods/%' AND size > 1024*100 GROUP BY md5 HAVING cnt < 10 ORDER BY cnt ASC;

SELECT * FROM scanned_files WHERE md5 IN (SELECT md5 as cnt from scanned_files WHERE path NOT LIKE '/etc/%' AND path NOT LIKE '/var/lib/docker/%' AND path NOT LIKE '/var/lib/kubelet/pods/%' AND size > 1024*100 GROUP BY md5 HAVING cnt < 10);
```

## Display file mode in octal format

```sql
SELECT path, size, CONV(mode, 10, 8) FROM scanned_files
```

## Find files that have setuid/setgid bits set

https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/stat.h#L17

```
#define S_IFMT  00170000
#define S_IFSOCK 0140000
#define S_IFLNK	 0120000
#define S_IFREG  0100000
#define S_IFBLK  0060000
#define S_IFDIR  0040000
#define S_IFCHR  0020000
#define S_IFIFO  0010000
#define S_ISUID  0004000
#define S_ISGID  0002000
#define S_ISVTX  0001000
```

```sql
SELECT m.id, m.hostname, f.id, f.path, CONV(f.mode, 10, 8) as mode, f.owner, f.group
FROM scanned_files f JOIN machines m ON f.machine_id = m.id
WHERE (mode & 2048) > 0 OR (mode & 1024) > 0;
```

## Find files in /usr/bin and group them by path and md5

```sql
SELECT path, md5, COUNT(*) as cnt FROM scanned_files WHERE path LIKE '/usr/bin/%' GROUP BY path, md5 ORDER BY path, cnt ASC
```

## Find unique files

```sql
SELECT m.id, m.hostname, f.id, f.path, CONV(f.mode, 10, 8) as mode, f.size, f.owner, f.group
FROM  scanned_files f JOIN machines m ON f.machine_id = m.id
GROUP BY f.md5
HAVING COUNT(*) = 1;

or

SELECT m.id, m.hostname, f.id, f.path, CONV(f.mode, 10, 8) as mode, f.size, f.owner, f.group
FROM (
    SELECT md5, machine_id, id
    FROM scanned_files
    GROUP BY md5, machine_id
    HAVING COUNT(*) = 1
) AS unique_md5_files
JOIN machines m ON unique_md5_files.machine_id = m.id
JOIN scanned_files f ON unique_md5_files.id = f.id;
```

## Find duplicate file submissions from same machine

```sql
SELECT path, md5, size, machine_id, COUNT(*) AS count
FROM scanned_files
GROUP BY path, md5, size, machine_id
HAVING count > 1;

SELECT path, md5, size, machine_id, GROUP_CONCAT(report_run_id) AS report_run_ids
FROM scanned_files
GROUP BY path, md5, size, machine_id
HAVING COUNT(*) > 1;

```
