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
FROM scanned_files f 
JOIN machines m ON f.machine_id = m.id
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


# Select latest submission from each machine for every path

```sql

SELECT t1.*
FROM scanned_files t1
WHERE t1.created_at = (
    SELECT MAX(t2.created_at)
    FROM scanned_files t2
    WHERE t2.path = t1.path AND t2.machine_id = t1.machine_id
)
GROUP BY t1.path, t1.machine_id, t1.created_at;

---

WITH latest_records AS (
    SELECT path, machine_id, MAX(created_at) as latest_created_at
    FROM scanned_files
    GROUP BY path, machine_id
)

SELECT t.*
FROM scanned_files t
JOIN latest_records lr
ON t.path = lr.path AND t.machine_id = lr.machine_id AND t.created_at = lr.latest_created_at;


---
SET autocommit=0;
INSERT INTO scanned_files_wip SELECT t1.*
FROM scanned_files t1
WHERE t1.created_at = (
    SELECT MAX(t2.created_at)
    FROM scanned_files t2
    WHERE t2.path = t1.path AND t2.machine_id = t1.machine_id
)
GROUP BY t1.path, t1.machine_id, t1.created_at;
COMMIT;
SET autocommit=1;

```


## Unique files by machine
    
```sql
SELECT m.hostname, m.fqdn, f.path, f.size, f.mtime, CONV(f.mode, 10, 8) as mode, f.owner as gid, f.`group` as gid 
FROM unique_files f 
JOIN machines m ON f.machine_id = m.id 
ORDER BY m.hostname, f.path
```

## Unqiue files over 100kB
```sql
SELECT `fqdn`, `path`, `size`,  CONV(`mode`, 10, 8) as mode, `mtime`, `owner`, `group`, `md5`, `sha1`, `sha256` 
FROM `unique_files` 
WHERE `size` > '100*1024'
ORDER BY fqdn;
```

## Select files that exists on multiple machines on same path, but with different content

```sql
SELECT `fqdn`, `path`, `size`, CONV(`mode`, 10, 8) as mode, `mtime`, `owner`, `group`, `md5`, `sha1`, `sha256` 
FROM unique_files 
WHERE id IN (
    SELECT id FROM unique_files 
    WHERE path NOT LIKE '/run/%' 
    AND path NOT LIKE '/var/log/journal/%' 
    AND path NOT LIKE '/opt/splunkforwarder/%' 
    GROUP BY path HAVING COUNT(*) > 1 
    )
ORDER BY fqdn
```


## Files that exists on multiple machines on same path, but with different permissions or owner or group

```sql
SELECT u1.fqdn as fqdn_1, u2.fqdn as fqdn_2, u1.`path`, CONV(u1.`mode`, 10, 8) as mode_1, CONV(u2.`mode`, 10, 8) as mode_2,  u1.`owner` as uid_1, u2.owner as uid_2, u1.`group` as gid_1, u2.`group` as gid_2, u1.`md5` as md5_1, u2.`md5` as md5_2 
FROM unique_files u1
JOIN unique_files u2
ON u1.path = u2.path
WHERE u1.mode != u2.mode
OR u1.owner != u2.owner
OR u1.`group` != u2.`group`
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

## Ideas for file hunting

```
SELECT m.id, m.hostname, f.id, f.path, CONV(f.mode, 10, 8) as mode, f.size, f.owner, f.group
FROM  scanned_files_fam f JOIN machines m ON f.machine_id = m.id
WHERE f.size > 1000 and m.hostname LIKE 'siem-%' AND f.path NOT LIKE '%wazuh%'AND f.path NOT LIKE '%splunk%' AND f.path NOT LIKE '/var/log/journal/%' AND f.path NOT LIKE '/var/ossec/%' AND ((f.mode & 1) > 0 OR ((f.mode >> 3) & 1) > 0 OR ((f.mode >> 6) & 1) > 0)
GROUP BY f.md5
HAVING COUNT(*) = 1
ORDER BY m.hostname, f.size
```


```sql

INSERT INTO known_files (path, sha1, sha256, md5, size, status)
SELECT scanned_files_wip.path, scanned_files_wip.sha1, scanned_files_wip.sha256, scanned_files_wip.md5, scanned_files_wip.size, 100 as status
FROM (
    SELECT sha1
    FROM scanned_files_wip
    GROUP BY sha1
    HAVING COUNT(*) >= 3
) AS subquery
JOIN scanned_files_wip ON subquery.sha1 = scanned_files_wip.sha1;


```