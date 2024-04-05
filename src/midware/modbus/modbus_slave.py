from pyModbusTCP.server import ModbusServer, DataBank
from pyModbusTCP.server import ModbusServer, DataHandler
from pyModbusTCP.constants import EXP_ILLEGAL_FUNCTION

from time import sleep
from random import uniform
import argparse
import csv
import struct

from pathlib import Path

from config import config

# list_of_objects = ["gwt-debug-tankItem1", "gwt-debug-tankItem2"]

objects = config.get("objects")

if __name__ == "__main__":

    # parse args
    #parser = argparse.ArgumentParser()
    #parser.add_argument('-H', '--host', type=str, default='127.0.0.1', help='Host (default: localhost)')
    #parser.add_argument('-p', '--port', type=int, default=5020, help='TCP port (default: 11502)')
    #args = parser.parse_args()
    
    # Create an instance of ModbusServer
    server = ModbusServer(host='0.0.0.0', port=5020, no_block=True)

    try:
        
        print("Server start...")
        server.start()
        
        while True:
    
            if objects:

                for index, item in enumerate(objects):
                    address = (index+1)*100
                    path = f'data/data_{index+1}.csv'
                    if Path(path).is_file():
                        with open(path, 'r') as file:
                            # Чтение значений как float
                            floats = list(map(lambda x: float(x), list(csv.reader(file, delimiter=',', quotechar='|', quoting=csv.QUOTE_MINIMAL))[0]))
                            print('from file' + str(floats))
                            for float_val in floats:
                                # Преобразование float в 16-битные регистры
                                reg_values = struct.unpack('HH', struct.pack('f', float_val))
                                # Запись в Modbus регистры
                                set_reg = server.data_bank.set_holding_registers(address, reg_values, srv_info=None)
                                address += 2  # Увеличиваем адрес, т.к. каждый float занимает 2 регистра

                    else:
                        print(f'file {path} not exist')    
                    sleep(1)

                sleep(10)          

    except Exception as exc:
        print(exc)
        print("Shutdown server ...")
        server.stop()
        print("Server is offline")  
