'use stirct'

var assert = require('assert')
var sqlQueryType = require('./')

var tests = [
  [null, ''],
  [undefined, ''],
  ['', ''],
  ['SELECT * FROM table_name', 'SELECT FROM table_name'],
  ['select * from table_name', 'SELECT FROM table_name'],
  ['  \n  SELECT  \n  * \n  FROM     \r\n     table_name  \n  \n', 'SELECT FROM table_name'],
  ['SELECT * FROM table_name; // foo', 'SELECT FROM table_name'],
  ['SELECT column1, column2 FROM table_name', 'SELECT FROM table_name'],
  ['SELECT DISTINCT column1, column2 FROM table_name', 'SELECT FROM table_name'],
  ['SELECT SUM(column_name) FROM table_name WHERE 1=1 GROUP BY column_name', 'SELECT FROM table_name'],
  ['SELECT COUNT(column_name) FROM table_name WHERE 1=1', 'SELECT FROM table_name'],
  ['CREATE TABLE table_name(column1 datatype, column2 datatype, column3 datatype, columnN datatype, PRIMARY KEY( one or more columns ))', 'CREATE TABLE table_name'],
  ['DROP TABLE table_name', 'DROP TABLE table_name'],
  ['CREATE INDEX index_name ON table_name ( column1, column2 )', 'CREATE INDEX index_name'],
  ['CREATE UNIQUE INDEX index_name ON table_name ( column1, column2 )', 'CREATE INDEX index_name'],
  ['ALTER TABLE table_name DROP INDEX index_name', 'ALTER TABLE table_name'],
  ['DESC table_name', 'DESC table_name'],
  ['TRUNCATE TABLE table_name', 'TRUNCATE TABLE table_name'],
  ['ALTER TABLE table_name ADD column_name datatype', 'ALTER TABLE table_name'],
  ['ALTER TABLE table_name RENAME TO new_table_name', 'ALTER TABLE table_name'],
  ['INSERT INTO table_name( column1, column2) VALUES ( value1, value2 )', 'INSERT INTO table_name'],
  ['UPDATE table_name SET column1 = 1, column2 = 2 WHERE 1=1', 'UPDATE table_name'],
  ['DELETE FROM table_name WHERE 1=1', 'DELETE FROM table_name'],
  ['CREATE DATABASE database_name', 'CREATE DATABASE database_name'],
  ['DROP DATABASE database_name', 'DROP DATABASE database_name'],
  ['USE database_name', 'USE database_name'],
  ['COMMIT', 'COMMIT'],
  ['ROLLBACK', 'ROLLBACK']
]

tests.forEach(function (test) {
  assert.strictEqual(sqlQueryType(test[0]), test[1])
})
