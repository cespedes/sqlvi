sqlvi is used to view and modfy a SQL database from a text editor.

In order to do it, `sqlvi` displays a list of rows and columns from
a SQL query as a text table, and launches an editor for the user.
After the editor finishes, it calculates if there are changes and, if so,
it updates the database accordingly.

In this state, it only works with PostgreSQL databases.

## Work in progress

This project is a work in progress.  Right now, it cannot do any changes in SQL:
it is only useful to display SQL results in a text editor.

## Details

When `sqlvi` executes, it does the following:
- Handle command-line arguments
- Read configuration file
- Connects to SQL database
- Executes SQL query
- Displays SQL result as an Org-mode text table
- Launches an editor for the user to see and modify the text table
- After the editor finishes, it finds the altered rows
- In a transaction, it inserts, updates, or deletes rows from the database

The SQL queries that `sqlvi` executes must be specified in the configuration file,
and the first column must be unique and not null (typically the primary key).

The header (first row) and the primary key (first column) should not be modified.
When a row is modified, it executes an *update*.
When a row is deleted, it executes a *delete*.
When the user inserts a new row, with the primary key empty, it executes an *insert*.

## Configuration file

`sqlvi` reads a YAML file (`$HOME/.sqlvi.yaml` by default) similar to this:

    editor: vim  # optional
    default: countries
    modes:
      countries:
        connect: postgres://user:pass@host/database
        select: SELECT id,country,capital,population FROM countries
        insert: INSERT INTO countries (country,capital,population) VALUES ($2,$3,$4)
        update: UPDATE countries SET country=$2,capital=$3,population=$4 WHERE id=$1
        delete: DELETE FROM countries WHERE id=$1
      products:
        connect: postgres://user:pass@host/database
        select: SELECT id,name,price FROM products
        insert: INSERT INTO products (name,price) VALUES ($2,$3)
        update: UPDATE products SET name=$2,price=$3 WHERE id=$1
        delete: DELETE FROM product WHERE id=$1

## Text table

The text table has this layout:

    |----+---------+---------+----------|
    | id | country | capital | language |
    |----+---------+---------+----------|
    | 1  | Spain   | Madrid  | Spanish  |
    | 2  | Germany | Berlin  | German   |
    | 3  | France  | Paris   | French   |
    |----+---------+---------+----------|

It should have these lines, and in this order:
- 0 or more lines that do not contain `|---` (which will be ignored).
- A separator line, beginning with `|---`.
- A line with list of column names, beginning and ending with `|` and separated by `|`.
- A separator line, beginning with `|---`.
- A line for each row, beginning and ending with `|`, and with values separated by `|`.
  The first field should be unique and non-empty.
- A separator line, beginning with `|---`.
- 0 or more lines after that one (ignored).
