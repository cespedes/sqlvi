sqlvi is used to view and modfy a SQL database from a text editor.

In order to do it, `sqlvi` performs a SQL query
and displays its result as a text file,
and launches an editor for the user.
After the editor finishes, it calculates if there are changes
and, if there are any,
it updates the database accordingly.

In this state, it only works with PostgreSQL databases.

## Details

When `sqlvi` executes, it does the following:
- Handle command-line arguments
- Read configuration file
- Connects to SQL database
- Executes SQL query
- Displays SQL result as a text-mode representation (text table, key-value fields...)
- Launches an editor for the user to see and modify the text table
- After the editor finishes, it finds the altered rows
- In a transaction, it inserts, updates, or deletes rows from the database

The SQL queries that `sqlvi` executes can be specified from command line or
in the configuration file.

The first field in the results must be unique and not null
(it is typically a primary key).
This field should not be modified from the editor.

When a row is modified, it executes an **update**.  
When a row is deleted, it executes a **delete**.  
When the user inserts a new row, with the first field empty, it executes an **insert**.

## Command-line arguments

The general usage of `sqlvi` is:

    sqlvi [options] [ <mode> ]

The `<mode>`, if given, specifies what set of options from the configuration file should
be applied.

The possible options are:

- `-debug`: enable debugging info.  Not very useful,
  unless you are debugging errors in the code.
- `-config <string>`: configuration file to read.
  If not specified, it reads `.sqlvirc` from the user's
  home directory.
- `-output <string>`: specify how to convert results to text, and vice versa.
  Possible values are
  `org` (Org-mode style text table, default),
  `ini` (Windows INI file)
  and
  `ldap` (paragraphs with "key: value" pairs).
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

`sqlvi` reads a YAML file (`$HOME/.sqlvirc` by default) similar to this:

    editor: vim
    connect: postgres://user:pass@host/database
    default: countries
    format: org
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

The configuration file is optional;
everything can be specified from command line,
which takes precedence over everything in this file.
The command-line arguments take precedenve over everything in the configuration file.

## Text output

There are several ways to display the query result in text, depending on
the value of the `output` configuration parameter.

### Org-mode table

If the `output` configuration parameter is `org`,
the results are rendered as a text table with this layout:

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

### INI file

If the `output` configuration parameter is `ini`,
the results are rendered as a INI-like file with this layout:

    [1]
    country = Spain
    capital = Madrid
    language = Spanish

    [2]
    country = Germany
    capital = Berlin
    language = German

    [3]
    country = France
    capital = Paris
    language = French

The order of the entries and the order of the fields in an entry is irrelevant.

An entry can be deleted or modified.

A new entry can be created using `[]` as section.

### YAML-like file

If the `output` configuration parameter is `ldap` or `yaml`,
the results are rendered as a list of paragraphs with `key: value` lines,
similar to YAML or the LDIF format using in LDAP:

    id: 1
    country: Spain
    capital: Madrid
    language: Spanish

    id: 2
    country: Germany
    capital: Berlin
    language: German

    id: 3
    country: France
    capital: Paris
    language: French

The order of the entries is irrelevant.
The first field in an entry should be the primary key (`id` in the example),
but the order of the rest of the fields is irrelevant.

An entry can be deleted or modified.

A new entry can be created using an empty value for the first key.

## Actions

After the editor finishes, `sqlvi` parses the text; if there are any errors,
it displays them and asks the user what to do (edit the file again, or quit).

If there are no errors, it finds out what lines, if any, have been modified.

If there are no modified lines, it shows `No changes.` and exits.

Otherwise, it displays a summary of changes, with the number of added,
updated and deleted lines, and asks the user what to do (edit the file again,
commit changes, or quit).

If the query to the database fails, it shows the error and asks the user again
what to do: edit the file again, or quit.
