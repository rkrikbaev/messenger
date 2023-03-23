from pyModbusTCP.server import ModbusServer, DataBank
from pyModbusTCP.server import ModbusServer, DataHandler
from pyModbusTCP.constants import EXP_ILLEGAL_FUNCTION

from time import sleep
from random import uniform
import argparse
import csv
import pdb

list_of_objects = ["gwt-debug-tankItem1", "gwt-debug-tankItem2"]
        
# class MyDataBank(DataBank):
#     """A custom ModbusServerDataBank for override get_holding_registers method."""

#     def __init__(self):
#         super().__init__()

#     def get_data(self):

#         for index, _ in enumerate(list_of_objects):
#             address = index*100
#             with open(f'data_{index+1}.csv', 'r') as file:
#                 word_list = csv.reader(file, delimiter=',', quotechar='|', quoting=csv.QUOTE_MINIMAL)
#                 print(word_list)
#             address = address + 1
#             return self.set_holding_registers(address, word_list)

if __name__ == "__main__":

    # parse args
    parser = argparse.ArgumentParser()
    parser.add_argument('-H', '--host', type=str, default='192.168.1.248', help='Host (default: localhost)')
    parser.add_argument('-p', '--port', type=int, default=5020, help='TCP port (default: 5020)')
    args = parser.parse_args()
    
    # Create an instance of ModbusServer
    server = ModbusServer(host=args.host, port=args.port, no_block=True)
    pdb.set_trace()
    try:
        print("Server start...")
        server.start()
        while True:
            for index, item in enumerate(list_of_objects):
                address = (index+1)*100
                with open(f'data_{index+1}.csv', 'r') as file:

                    word_list = list(map(lambda x: int(x), list(csv.reader(file, delimiter=',', quotechar='|', quoting=csv.QUOTE_MINIMAL))[0]))
                    print('from file' + str(list(word_list)))

                set_reg = server.data_bank.set_holding_registers(address, word_list, srv_info=None)
                print('written to h_reg: ' + str(set_reg))
                # res = server.data_bank.get_holding_registers(0,5)
                # print('red from h_reg: ' + str(res))
                sleep(1) 
            sleep(10)            

    except Exception as exc:
        print(exc)
        print("Shutdown server ...")
        server.stop()
        print("Server is offline")  