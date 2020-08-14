import json
import mysql.connector
import os

mydb = mysql.connector.connect(
  host=os.environ.get("MYSQL_HOST"),
  user=os.environ.get("MYSQL_USER"),
  password=os.environ.get("MYSQL_PASSWORD"),
  database="quilt"
)

c = mydb.cursor()

columns = [
	"id",
	"peer",
	"is_direct",
	"submitted",
	"duplicate",
	"otherreject",
	"underpriced",
	"t",
]

create_table = """
	CREATE TABLE if not EXISTS wire_transactions (
		id TEXT,
		peer VARCHAR(256),
		is_direct BOOLEAN,
		submitted INTEGER,
		duplicate INTEGER,
		otherreject INTEGER,
		underpriced INTEGER,
		ts DATETIME
	);
"""

insert_query = """
	INSERT into wire_transactions(
		id,
		peer,
		is_direct,
		submitted,
		duplicate,
		otherreject,
		underpriced,
		ts
	) values (?, ?, ?, ?, ?, ?, ?, ?)
"""

create_index = """
	CREATE INDEX ts
	ON wire_transactions (ts);
"""

c.execute(create_table)

data = []
with open('./wire_events.json') as f:
	for row in f:
		line = json.loads(row)
		data.append(line)

ordered_columns = []
for i, entry in enumerate(data):
	values = []
	for j, column in enumerate(columns):
		if column == "id":
			values.append(i)
		else:
			values.append(entry[column])
	ordered_columns.append(tuple(values))

c.executemany(insert_query, ordered_columns)
c.execute(create_index)

mydb.commit()
