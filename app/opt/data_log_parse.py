import csv

object_name = input("name: ").lower()

# Открытие исходного файла CSV
with open(f'data_log - {object_name}.csv', 'r') as file:
    reader = csv.reader(file)
    data = list(reader)
    print(data)

# Получение названий столбцов и удаление из массива данных
headers = data.pop(0)
print(headers)

headers_enable = ['dev1_density', 'dev1_volume', 'dev1_temperature', 'dev1_massflowbegin', 'dev1_massflowend', 'dev1_mass', 'dev1_masstotalizer', 'dev2_density', 'dev2_volume', 'dev2_temperature', 'dev2_tankLevel', 'dev2_mass', 'dev3_density', 'dev3_volume', 'dev3_temperature', 'dev3_tankLevel', 'dev3_mass']

# Разделение данных и сохранение в отдельные файлы CSV
for i, row in enumerate(data):
    
    data_to_isun = [ x for index, x in enumerate(row) if headers[index] in headers_enable ]
    #print(row)
    print(data_to_isun)
    
    date_value = row[2]
    
    date_value = date_value.replace("/",".")
    print(date_value)
    new_data = [[cell[0], cell[1]] for cell in zip(headers_enable, data_to_isun)]
    new_filename = f'log_{date_value}.csv'
    try:
        with open(new_filename, 'w', newline='') as file:
            writer = csv.writer(file)
            writer.writerows(new_data)
        
        print(f'Файл {new_filename} сохранен.')
    except:
        pass
