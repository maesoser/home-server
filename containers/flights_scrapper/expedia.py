from time import sleep, strftime
import time
from random import randint
import pandas as pd
from selenium.webdriver.common.keys import Keys
import smtplib
from email.mime.multipart import MIMEMultipart
from selenium import webdriver
import datetime

browser_options = webdriver.ChromeOptions()
browser_options.add_argument("--no-sandbox")
browser = webdriver.Remote(
    command_executor='http://192.168.0.10:4444/wd/hub',
    options=browser_options
)

#Setting ticket types paths
return_ticket = "//label[@id='flight-type-roundtrip-label-hp-flight']"
one_way_ticket = "//label[@id='flight-type-one-way-label-hp-flight']"
multi_ticket = "//label[@id='flight-type-multi-dest-label-hp-flight']"

def ticket_chooser(ticket):
    try:
        ticket_type = browser.find_element_by_xpath(ticket)
        ticket_type.click()
    except Exception as e:
        pass
def dep_country_chooser(dep_country):
    fly_from = browser.find_element_by_xpath("//input[@id='flight-origin-hp-flight']")
    time.sleep(1)
    fly_from.clear()
    time.sleep(1.5)
    fly_from.send_keys('  ' + dep_country)
    time.sleep(1.5)
    first_item = browser.find_element_by_xpath("//a[@id='aria-option-0']")
    time.sleep(1.5)
    first_item.click()

def arrival_country_chooser(arrival_country):
    fly_to = browser.find_element_by_xpath("//input[@id='flight-destination-hp-flight']")
    time.sleep(1)
    fly_to.clear()
    time.sleep(1.5)
    fly_to.send_keys('  ' + arrival_country)
    time.sleep(1.5)
    first_item = browser.find_element_by_xpath("//a[@id='aria-option-0']")
    time.sleep(1.5)
    first_item.click()

def dep_date_chooser(month, day, year):
    dep_date_button = browser.find_element_by_xpath("//input[@id='flight-departing-hp-flight']")
    dep_date_button.clear()
    dep_date_button.send_keys(month + '/' + day + '/' + year)

def return_date_chooser(month, day, year):
    return_date_button = browser.find_element_by_xpath("//input[@id='flight-returning-hp-flight']")
    for i in range(11):
        return_date_button.send_keys(Keys.BACKSPACE)
    return_date_button.send_keys(month + '/' + day + '/' + year)

def search():
    search = browser.find_element_by_xpath("//button[@class='btn-primary btn-action gcw-submit']")
    search.click()
    time.sleep(15)
    print('Results ready!')

df = pd.DataFrame()
def compile_data():
    global df
    global dep_times_list
    global arr_times_list
    global airlines_list
    global price_list
    global durations_list
    global stops_list
    global layovers_list
    #departure times
    dep_times = browser.find_elements_by_xpath("//span[@data-test-id='departure-time']")
    dep_times_list = [value.text for value in dep_times]
    #arrival times
    arr_times = browser.find_elements_by_xpath("//span[@data-test-id='arrival-time']")
    arr_times_list = [value.text for value in arr_times]
    #airline name
    airlines = browser.find_elements_by_xpath("//span[@data-test-id='airline-name']")
    airlines_list = [value.text for value in airlines]
    #prices
    prices = browser.find_elements_by_xpath("//span[@data-test-id='listing-price-dollars']")
    print(prices[0].text)
    price_list = [value.text for value in prices]
    #durations
    durations = browser.find_elements_by_xpath("//span[@data-test-id='duration']")
    durations_list = [value.text for value in durations]
    #stops
    stops = browser.find_elements_by_xpath("//span[@class='number-stops']")
    stops_list = [value.text for value in stops]
    #layovers
    layovers = browser.find_elements_by_xpath("//span[@data-test-id='layover-airport-stops']")
    layovers_list = [value.text for value in layovers]
    now = datetime.datetime.now()
    current_date = (str(now.year) + '-' + str(now.month) + '-' + str(now.day))
    current_time = (str(now.hour) + ':' + str(now.minute))
    current_price = 'price' + '(' + current_date + '---' + current_time + ')'
    for i in range(len(dep_times_list)):
        try:
            df.loc[i, 'departure_time'] = dep_times_list[i]
        except Exception as e:
            pass
        try:
            df.loc[i, 'arrival_time'] = arr_times_list[i]
        except Exception as e:
            pass
        try:
            df.loc[i, 'airline'] = airlines_list[i]
        except Exception as e:
            pass
        try:
            df.loc[i, 'duration'] = durations_list[i]
        except Exception as e:
            pass
        try:
            df.loc[i, 'stops'] = stops_list[i]
        except Exception as e:
            pass
        try:
            df.loc[i, 'layovers'] = layovers_list[i]
        except Exception as e:
            pass
        try:
            df.loc[i, str(current_price)] = price_list[i]
        except Exception as e:
            pass
    print('Excel Sheet Created!')

#Create message template for email
def create_msg():
    global msg
    msg = '\nCurrent Cheapest flight:\n\nDeparture time: {}\nArrival time: {}\nAirline: {}\nFlight duration: {}\nNo. of stops: {}\nPrice: {}\n'.format(cheapest_dep_time,
                       cheapest_arrival_time,
                       cheapest_airline,
                       cheapest_duration,
                       cheapest_stops,
                       cheapest_price)

for i in range(1):    
    print("Scrapping")
    link = 'https://www.expedia.com/'
    browser.get(link)
    time.sleep(5)
    #choose flights only
    flights_only = browser.find_element_by_xpath("//button[@id='tab-flight-tab-hp']")
    flights_only.click()
    ticket_chooser(return_ticket)
    dep_country_chooser('London')
    arrival_country_chooser('Madrid')
    today = datetime.date.today()
    return_date = today + datetime.timedelta(days=7)
    dep_date_chooser(today.strftime("%m"), today.strftime("%d"), str(today.year))
    return_date_chooser(return_date.strftime("%m"), return_date.strftime("%d"), str(return_date.year))
    search()
    compile_data()
    #print(df)
    #save values for email
    current_values = df.iloc[0]
    cheapest_dep_time = current_values[0]
    cheapest_arrival_time = current_values[1]
    cheapest_airline = current_values[2]
    cheapest_duration = current_values[3]
    cheapest_stops = current_values[4]
    cheapest_price = current_values[-1]
    print('run {} completed!'.format(i))
    create_msg()
    df.to_csv('flights.csv',sep='\t')
    time.sleep(3600)


