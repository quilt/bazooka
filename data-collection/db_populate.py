#!/Library/Frameworks/Python.framework/Versions/3.8/bin/python3

import sqlite3
import json
import uuid

conn = sqlite3.connect('aa-collection.db')

columns = [
	"id",
	"run_id",
	"is_aa",
	"hash",
	"data",
	"gas_price",
	"recipient",
	"gas",
	"gas_used",
	"start_time",
	"end_time",
	"duration",
	"success",
]


create_table = """
	CREATE TABLE if not EXISTS validations (
		id TEXT,
		run_id INTEGER,
		is_aa BOOLEAN,
		hash VARCHAR(256),
		data TEXT,
		gas_price TEXT,
		recipient TEXT,
		gas TEXT,
		gas_used TEXT,
		start_time INTEGER,
		end_time INTEGER,
		duration INTEGER,
		success BOOLEAN
	);
"""

run_id_query = """
	SELECT
		IFNULL(MAX(run_id), -1)
	FROM validations;
"""

c = conn.cursor()
c.execute(create_table)

run_id = c.execute(run_id_query).fetchone()

data = []
with open('./data.json') as f:
	for row in f:
		line = json.loads(row)
		if line["run_id"] > run_id[0]:
			data.append(line)

ordered_columns = []

for i, entry in enumerate(data):
	values = []
	for j, column in enumerate(columns):
		if column == "id":
			values.append(str(uuid.uuid4()))
		else:
			values.append(entry[column])
	ordered_columns.append(tuple(values))

insert_query = """
	INSERT into validations(
		id,
		run_id,
		is_aa,
		hash,
		data,
		gas_price,
		recipient,
		gas,
		gas_used,
		start_time,
		end_time,
		duration,
		success
	) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
"""

c.executemany(insert_query, ordered_columns)

conn.commit()
c.close()








