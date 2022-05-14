sqlvi is used to view and modfy a SQL database from a text editor.

In order to do it, `sqlvi` displays a list of rows and columns from
a SQL query as a text table, and launches an editor for the user.
After the editor finishes, it calculates if there are changes and, if so,
it updates the database accordingly.

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
- In a transaction, it inserts, updates, or drops rows from the database

The SQL queries that `sqlvi` executes must be specified in the configuration file,
and the first row must be a primary key, not null and different from all the other lines.

## Conficuration file

`sqlvi` reads a YAML file (`$HOME/.sqlvi.yaml` by default) similar to this:

    editor: vim  # optional
    default: countries
    modes:
      countries:
        connect: postgres://user:pass@host/database
        select: SELECT id,country,capital,population FROM countries
        insert: INSERT INTO countries (country,capital,population) VALUES ($2,$3,$4)
        update: UPDATE countries SET country=$2,capital=$3,population=$4 WHERE id=$1
        drop: DROP countries WHERE id=$1
      products:
        connect: postgres://user:pass@host/database
        select: SELECT id,name,price FROM products
        insert: INSERT INTO products (name,price) VALUES ($2,$3)
        update: UPDATE products SET name=$2,price=$3 WHERE id=$1
        drop: DROP product WHERE id=$1

## Text table

The text table has this layout:

    |----+---------+---------+------------|
    | id | country | capital | population |
    |----+---------+---------+------------|
    | 1  | Spain   | Madrid  | 47450795   |
    | 2  | Germany | Berlin  | 83190556   |
    |----+---------+---------+------------|

It should have:
- A separator line, beginning with `|--`.
- A line with list of column names, beginning and ending with `|` and separated by `|`.
- A separator line, beginning with `|--`.
- A line for each row, beginning and ending with `|`, and with values separated by `|`.
- A separator line, beginning with `|--`.
