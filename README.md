sqlvi is used to view and modfy a SQL database from a text editor.

In order to do it, `sqlvi` displays a list of rows and columns from
a SQL query as a text table, and launches an editor for the user.
After the editor finishes, it calculates if there are changes and, if so,
it updates the database accordingly.

In this state, it only works with PostgreSQL databases.

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

The SQL queries that `sqlvi` executes can be specified from command line or
in the configuration file.

The first column in the results must be unique and not null
(it is typically a primary key).

The first row (header) and the first column should not be modified from the editor.
When a row is modified, it executes an **update**.
When a row is deleted, it executes a **delete**.
When the user inserts a new row, with the first column empty, it executes an **insert**.

## Command-line arguments

The general usage of `sqlvi` is:

    sqlvi [options] [ <mode> ]

The `<mode>`, if given, specifies what set of options from the configuration file should
be applied.

The possible options are:

- `-debug`: enable debugging info.  Not very useful,
  unless you are debugging errors in the code.
- `-config <string>`: configuration file to read.
  If not specified, it reads `.sqlvi.yaml` from the user's
  home directory.
- `-connect <string>`: connection string to the database.
  Right now it can only be user with PostgreSQL databases.
- `-editor <string>`: editor to use.  If not specified
  (and not found in the configuration file), it tries `$VISUAL`,
  `$EDITOR`, `editor`, `vim` and `vi`, in that order.
- `-select <string>`: Database query to execute on startup and
  display its result.
- `-insert <string>`: Database query to add a new row.
- `-update <string>`: Database query to update a row.
- `-delete <string>`: Database query to delete a row.

In the `insert`, `update` and `delete` entries there can be arguments
specified as `$1`, `$2`, `$3`... which will be replaced by those
columns in the modified row.  It is not needed to specify all the columns, or to
have them in the same order they appear.  They can also be repeated.

## Configuration file

`sqlvi` reads a YAML file (`$HOME/.sqlvi.yaml` by default) similar to this:

    editor: vim  # optional
    default: countries
    modes:
      countries:
        connect: postgres://user:pass@host/database
        select: SELECT id,country,capital,language FROM countries
        insert: INSERT INTO countries (country,capital,language) VALUES ($2,$3,$4)
        update: UPDATE countries SET country=$2,capital=$3,language=$4 WHERE id=$1
        delete: DELETE FROM countries WHERE id=$1
      products:
        connect: postgres://user:pass@host/database
        select: SELECT id,name,price FROM products
        insert: INSERT INTO products (name,price) VALUES ($2,$3)
        update: UPDATE products SET name=$2,price=$3 WHERE id=$1
        delete: DELETE FROM product WHERE id=$1

The command-line arguments take precedenve over everything in the configuration file.

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

## Actions

After the editor finishes, `sqlvi` parses the edited table; if there are any errors,
it displays them and asks the user what to do (edit the file again, or quit).

If there are no errors, it finds out what lines, if any, have been modified.

If there are no modified lines, it shows `No changes.` and exits.

Otherwise, it displays a summary of changes, with the number of added,
updated and deleted lines, and asks the user what to do (edit the file again,
commit changes, or quit).

If the query to the database fails, it shows the error and asks the user again
what to do: edit the file again, or quit.
