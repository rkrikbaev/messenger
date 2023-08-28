# Скрипт из определеной директории берет CSV файл 
# содержащий показания массы в двух емкостях за каждый день
# На основании изменеия массы вычисляется откачка 
# топлива на отгрузку и записывается обрабтно в файл
# Поля в файле: 'Date', 
#               'dev1_density', 'dev1_volume', 'dev1_temperature', 'dev1_massflowbegin', 'dev1_massflowend', 'dev1_mass', 'dev1_masstotalizer', 
#               'dev2_density', 'dev2_volume', 'dev2_temperature', 'dev2_level', 'dev2_mass', 
#               'dev3_density', 'dev3_volume', 'dev3_temperature', 'dev3_level', 'dev3_mass'
# dev1, dev2 и dev3 - отгрузка, хранение в емкости 1 и 2

import csv
import os
from datetime import datetime
import re

def calculate_flow(mass_end, mass_begin):
    mass_upload = mass_end - mass_begin
    if mass_upload >= 0:
        return 0, mass_begin, mass_end
    else:
        mass_end = mass_begin
        return abs(mass_upload), mass_begin, mass_end

def process_number_density(number):

    # Проверка количества знаков в целой части
    if int(number) > 0:
        # Целая часть содержит как минимум 1 знак
        int(number)
        return number
    else:
        # Целая часть не содержит знаков, умножаем на 100
        return number * 1000

def find_files_with_extension(folder, extension):
    files_with_extension = []
    for filename in os.listdir(folder):
        if filename.endswith(extension):
            files_with_extension.append(os.path.join(folder, filename))
    return files_with_extension

def change_date(date_str):
    # Преобразуем строку даты в объект datetime
    date = datetime.strptime(date_str, "%Y.%m.%d")
    # Форматируем обратно в нужный формат
    formatted_date = date.strftime("%Y/%m/%d")
    print("Преобразованная дата:", formatted_date)
    return formatted_date

def main(folder_path, out_folder):
    # Получение списка файлов в папке
    files = find_files_with_extension(folder_path, '.csv')
    for file_path in files:
        array_data = []
        with open(file_path, 'r', newline='\n') as csvfile:
            reader = csv.reader(csvfile, delimiter=';', quoting=csv.QUOTE_NONE)         
            # Убираем заголовки (названий колонок)
            _ = next(reader)
            array_data.append(['Date', 
                            'dev1_density','dev1_volume','dev1_temperature','dev1_massflowbegin','dev1_massflowend','dev1_mass',
                            'dev2_density','dev2_volume','dev2_temperature','dev2_level','dev2_mass',
                            'dev3_density','dev3_volume','dev3_temperature','dev3_level','dev3_mass'])
            # set zero for mass counter
            dev2_mass_upload, dev2_mass_begin, dev2_mass_end = 0, 0, 0
            dev3_mass_upload, dev3_mass_begin, dev3_mass_end = 0, 0, 0
            mass_begin, mass_end, mass, totalizer = 0, 0, 0, 0           
            for index, row in enumerate(reader):
                if index == 0:
                    first_date = str(row[0])                    
                    year, month, _ = first_date.split('.')
                    file_date = f'{year}.{month}'

                date = change_date(row.pop(0))
                _ = row.pop(0)
                _ = row.pop(4)
                num_row = [float(item.replace(',', '.')) for item in row]

                dev1_density = 0
                dev1_temperature = 0

                dev2_density = process_number_density(num_row[0])
                dev2_volume = num_row[2]
                dev2_temperature = 0
                dev2_level = num_row[1]
                dev2_mass = num_row[3]

                dev3_density = process_number_density(num_row[4])
                dev3_volume = num_row[6]
                dev3_temperature = 0
                dev3_level = num_row[5]
                dev3_mass = num_row[7]

                dev2_mass_end = dev2_mass
                dev3_mass_end = dev3_mass

                # calculate parameters for massflow
                dev2_mass_upload, _, _ = calculate_flow(mass_end=dev2_mass_end, mass_begin=dev2_mass_begin)
                dev3_mass_upload, _, _ = calculate_flow(mass_end=dev3_mass_end, mass_begin=dev3_mass_begin)

                # parameters take from tank
                if dev2_mass_upload > 0:
                    dev1_density = dev2_density
                    dev1_temperature = dev2_temperature  
                elif dev3_mass_upload > 0:
                    dev1_density = dev2_density
                    dev1_temperature = dev2_temperature  
                elif (dev2_mass_upload > 0) and (dev3_mass_upload > 0):
                    dev1_density = ( dev2_density + dev3_density ) / 2
                    dev1_temperature = ( dev2_temperature + dev3_temperature ) / 2
                else:
                    dev1_density = 0
                    dev1_temperature = 0                 

                dev1_volume = 0
                mass = dev2_mass_upload + dev3_mass_upload
                mass_end = mass_end + mass

                array_data.append([date, 
                                dev1_density, dev1_volume, dev1_temperature, mass_begin, mass_end, mass,
                                dev2_density, dev2_volume, dev2_temperature, dev2_level, dev2_mass,
                                dev3_density, dev3_volume, dev3_temperature, dev3_level, dev3_mass])

                mass_begin = mass_end
                dev2_mass_begin = dev2_mass_end
                dev3_mass_begin = dev3_mass_end

            new_file_name = [(f'data_{file_date}.csv', array_data)]
            
            for file_name, data in new_file_name:
                print(out_folder+file_name)
                with open(out_folder+file_name, 'w') as csvfile:
                    writer = csv.writer(csvfile,)
                    writer.writerows(data)


if __name__ == '__main__':

    subfolders = ['aman', 'ayraqty', 'zhar'] #input('out folder name: ')
    years = ['2021', '2022', '2023'] #input('enter year: ')

    for subfolder in subfolders:
        for year_value in years:
            try:
                folder_path = f'data_past/in/{subfolder}/tank/{year_value}'
                out_folder = f'data_past/out/{subfolder}/'
                # Проверяем, существует ли папка
                if not os.path.exists(out_folder):
                    # Создаем папку, если она не существует
                    os.makedirs(out_folder)
                    print("Папка создана.")
                else:
                    pass
                    # print("Папка уже существует.")
                main(folder_path, out_folder)
            except:
                pass

