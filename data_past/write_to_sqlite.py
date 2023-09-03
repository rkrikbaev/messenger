import os
import sqlite3
from datetime import datetime
import csv

# Имя базы данных SQLite
database_name = "mydatabase.db"

# Подключение к базе данных
conn = sqlite3.connect(database_name)
cursor = conn.cursor()

def find_files_with_extension(folder, extension):
    files_with_extension = []
    for filename in os.listdir(folder):
        if filename.endswith(extension):
            files_with_extension.append(os.path.join(folder, filename))
    return files_with_extension

# The function return list of the folders that do not have any subfolders.
def list_folders_recursively(folder_path):
    list_of_paths = []
    list_of_names = []
    for root, folders, files in os.walk(folder_path):
        if not folders:  # Check if the current folder has no subdirectories
            print(root)
            list_of_paths.append(root)
            list_of_names.append(os.path.basename(root))
    return list_of_paths, list_of_names

# Check if the table exists
def table_exists(table_name):
    cursor.execute("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", (table_name,))
    return cursor.fetchone()[0] == 1

# Check if the Date value already exists in the table
def date_exists(table_name, date_value):
    # SELECT query
    query = f"SELECT date FROM {table_name}"
    cursor.execute(query)
    # Fetch all the rows
    results = cursor.fetchall()
    # Convert to a list of strings
    date_list = [row[0] for row in results]
    if date_value in date_list:
        return True
    else:
        return False

# Write to db
def write_to_db(table_name, data):
    if not table_exists(table_name):
        # Create a table if it doesn't exist
        create_table_sql = f'''
            CREATE TABLE {table_name} (
                Date DATETIME PRIMARY KEY,
                dev1_density REAL,
                dev1_volume REAL,
                dev1_temperature REAL,
                dev1_massflowbegin REAL,
                dev1_massflowend REAL,
                dev1_mass REAL,
                dev1_masstotalizer REAL,
                dev2_density REAL,
                dev2_volume REAL,
                dev2_temperature REAL,
                dev2_level REAL,
                dev2_mass REAL,
                dev3_density REAL,
                dev3_volume REAL,
                dev3_temperature REAL,
                dev3_level REAL,
                dev3_mass REAL
            )
        '''
        cursor.execute(create_table_sql)

    # Convert the date string to a Datetime object
    date_value = data[0].split(',')[0]
    values_to_write = data[0].split(',')
    columns = "'Date', 'dev1_density', 'dev1_volume', 'dev1_temperature', 'dev1_massflowbegin', 'dev1_massflowend', 'dev1_mass', 'dev1_masstotalizer', 'dev2_density', 'dev2_volume', 'dev2_temperature', 'dev2_level', 'dev2_mass', 'dev3_density', 'dev3_volume', 'dev3_temperature', 'dev3_level', 'dev3_mass'"
    if not date_exists(table_name, date_value):
        # Insert data into the table if the Date value doesn't exist
        # Add single quotes around the date value in the INSERT statement
        insert_data_sql = f"INSERT INTO {table_name} ({columns}) VALUES ({','.join(['?'] * len(values_to_write))})"
        cursor.execute(insert_data_sql, values_to_write)
        print("Data inserted into the database.")
    else:
        print("Data with the same Date value already exists in the database.")
        # Сохранение изменений и закрытие соединения с базой данных

def main(folder_path):

    # Получение списка файлов в папке
    files = os.listdir(folder_path)

    root_folder_path = folder_path
    paths, names = list_folders_recursively(root_folder_path)

    # Проход по каждому файлу
    for path_to_csv, table_name in zip(paths, names):
        files = sorted(find_files_with_extension(path_to_csv, extension='.csv'))
        for file in files:
            with open(file, 'r') as csvfile:
                reader = csv.reader(csvfile, delimiter=';', quoting=csv.QUOTE_NONE)
                column_names = next(reader)
                print(column_names[0].split(',')) 
                for index, row in enumerate(reader):
                    # print(row[0])    
                    write_to_db(table_name, data=row)
                conn.commit()
    conn.close()

if __name__ == "__main__":

    # Путь к папке с файлами
    folder_path = "./data_past/out/"
    main(folder_path)
            