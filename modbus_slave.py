from pyModbusTCP.server import ModbusServer, DataBank
from pyModbusTCP.server import ModbusServer, DataHandler
from pyModbusTCP.constants import EXP_ILLEGAL_FUNCTION

from time import sleep
from random import uniform
import argparse
import csv

list_of_objects = ["gwt-debug-tankItem1", "gwt-debug-tankItem2"]
        
class MyDataBank(DataBank):
    """A custom ModbusServerDataBank for override get_holding_registers method."""

    def __init__(self):
        super().__init__()

    def get_data(self):

        for index, item in enumerate(list_of_objects):

            with open(f'data_{index+1}.csv', 'r') as file:
                word_list = csv.reader(file, delimiter=',', quotechar='|', quoting=csv.QUOTE_MINIMAL)
                address = index*100
                print(word_list)

            return self.set_holding_registers(address, word_list)

if __name__ == "__main__":

    # parse args
    parser = argparse.ArgumentParser()
    parser.add_argument('-H', '--host', type=str, default='localhost', help='Host (default: localhost)')
    parser.add_argument('-p', '--port', type=int, default=11502, help='TCP port (default: 11502)')
    args = parser.parse_args()
    
    # Create an instance of ModbusServer
    server = ModbusServer(host=args.host, port=args.port, data_bank=MyDataBank())
    try:
        print("Server start...")
        server.start()

    except Exception as exc:
        print(exc)
        print("Shutdown server ...")
        server.stop()
        print("Server is offline")  