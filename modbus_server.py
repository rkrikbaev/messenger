from pyModbusTCP.server import ModbusServer, DataBank
from time import sleep
from random import uniform

from bs4 import BeautifulSoup
import requests

from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options

url_login = "http://127.0.0.1:5500/Veeder-Root%20Web%20Interface_Login_files/Veeder-Root%20Web%20Interface_Login.html"
url_data = 'http://127.0.0.1:5500/Veeder-Root%20Web%20Interface/Veeder-Root%20Web%20Interface.html'

def set_chrome_options() -> None:
    """Sets chrome options for Selenium.
    Chrome options for headless browser is enabled.
    """
    chrome_options = Options()
    chrome_options.add_argument("--headless")
    # chrome_options.add_argument("--no-sandbox")
    # chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_prefs = {}
    chrome_options.experimental_options["prefs"] = chrome_prefs
    chrome_prefs["profile.default_content_settings"] = {"images": 2}
    
    return chrome_options

def login(url):
    
    driver.get(url)

    # To catch <input type="text" id="gwt-debug-userPasswordTextBox" />
    password = driver.find_element(By.ID, "gwt-debug-userPasswordTextBox")
    # To catch <input type="text" id="gwt-debug-userNameTextBox" />
    login = driver.find_element(By.ID, "gwt-debug-userNameTextBox")

    password.send_keys("admin")
    login.send_keys("administrator")

    driver.find_element(By.ID, "gwt-debug-signInButton").click()

    return

def get_data(item, url):

    try:
        driver.get(url)
        r = requests.get().content
    except:
        pass
    
    content= BeautifulSoup(r, 'html.parser')

    find_text= content.find("div", {"id": item}).text

    list_of = find_text.split('  ')

    def convert_to_int(element):

        if isinstance(element, str):
        
            try:
                element = float(element)
            except ValueError:
                element = element.strip()
            
            return element

    def convert_to_dict(a):

        it = iter(a)
        res_dct = dict(zip(it, it))
        return res_dct

    list_of_values = [convert_to_int(x) for x in list_of if x!='' and x.strip() != 'View']

    label = list_of_values.pop(0)

    values = convert_to_dict(list_of_values)

    print(label)
    print(values)

    return label, values


if __name__ == "__main__":

    driver = webdriver.Chrome('chromedriver')

    login(url_login)
    
    # Create an instance of ModbusServer
    server = ModbusServer("127.0.0.1", 12345, no_block=True)
    list_of_objects = ["gwt-debug-tankItem1", "gwt-debug-tankItem2"]

    try:
        print("Start server...")
        server.start()
        print("Server is online")
        state = [0]
        while True:
            for item in list_of_objects:
                
                d = get_data(item, url_data)
                print(get_data(d))

            DataBank.set_words(0, [int(uniform(0, 100))])
            if state != DataBank.get_words(1):
                state = DataBank.get_words(1)
                print("Value of Register 1 has changed to " +str(state))
            sleep(2)

    except:
        print("Shutdown server ...")
        server.stop()
        print("Server is offline")    


    driver.close()