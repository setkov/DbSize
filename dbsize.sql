USE master

SELECT	Database_Name = d.name,
		Total_size_mb = CAST(SUM(size) * 8. / 1024 AS DECIMAL(8,2)),
		Row_size_mb = CAST(SUM(CASE WHEN type_desc = 'ROWS' THEN size END) * 8. / 1024 AS DECIMAL(8,2)),
		Log_size_mb = CAST(SUM(CASE WHEN type_desc = 'LOG' THEN size END) * 8. / 1024 AS DECIMAL(8,2)),
		Created = FORMAT(d.create_date, 'yyyy-MM-dd'),
		Owner = SUSER_SNAME(d.owner_sid),
		State = CASE 
					WHEN d.state = 0 THEN 'ONLINE'
					WHEN d.state = 1 THEN 'RESTORING'
					WHEN d.state = 2 THEN 'RECOVERING'
					WHEN d.state = 3 THEN 'RECOVERY_PENDING'
					WHEN d.state = 4 THEN 'SUSPECT'
					WHEN d.state = 5 THEN 'EMERGENCY'
					WHEN d.state = 6 THEN 'OFFLINE'
					WHEN d.state = 7 THEN 'COPYING'
					WHEN d.state = 10 THEN 'OFFLINE_SECONDARY'
					ELSE 'UNKNOWN'
				END,
		MS_Description = CAST('' AS SQL_VARIANT)
INTO	#tables
FROM	sys.master_files f WITH(NOWAIT)
	INNER JOIN sys.databases d
		ON d.database_id = f.database_id
GROUP BY d.name, d.create_date, d.owner_sid,
		CASE 
			WHEN d.state = 0 THEN 'ONLINE'
			WHEN d.state = 1 THEN 'RESTORING'
			WHEN d.state = 2 THEN 'RECOVERING'
			WHEN d.state = 3 THEN 'RECOVERY_PENDING'
			WHEN d.state = 4 THEN 'SUSPECT'
			WHEN d.state = 5 THEN 'EMERGENCY'
			WHEN d.state = 6 THEN 'OFFLINE'
			WHEN d.state = 7 THEN 'COPYING'
			WHEN d.state = 10 THEN 'OFFLINE_SECONDARY'
			ELSE 'UNKNOWN'
		END

DECLARE	@value SQL_VARIANT,
		@property NVARCHAR(20) = 'MS_Description',
		@sql NVARCHAR(200),
		@database sysname

DECLARE tables CURSOR FAST_FORWARD FOR 
SELECT [Database_Name] FROM #tables WHERE [State] = 'ONLINE'
  
OPEN tables 
FETCH NEXT FROM tables INTO @database
WHILE @@FETCH_STATUS = 0  
BEGIN  
	SELECT	@value = NULL,
			@sql = CONCAT('use [', @database, ']; SELECT @value=value FROM sys.extended_properties WHERE class = 0 AND name = @property')
	EXECUTE sp_executesql @sql, N'@property VARCHAR(20), @value SQL_VARIANT OUTPUT', @property= @property, @value = @value OUTPUT
--	SELECT @database, @sql, @value

	UPDATE #tables SET MS_Description = ISNULL(@value, '') WHERE Database_Name = @database

	FETCH NEXT FROM tables INTO @database
END  
CLOSE tables
DEALLOCATE tables

SELECT * FROM #tables ORDER BY Total_size_mb DESC
DROP TABLE #tables