import requests
import pymysql
from bs4 import BeautifulSoup

# Database connection details
db_host = 'localhost'
db_user = 'root'
db_password = ''
db_name = 'kajaki'
table_name = 'wyscigi'

def fetch_data(url):
    response = requests.get(url)
    if response.status_code == 200:
        return response.content.decode('iso-8859-2')  # Decode using ISO-8859-2
    else:
        print(f"Failed to fetch data: Status code {response.status_code}")
        return None

def parse_data(html_data):
    soup = BeautifulSoup(html_data, 'html.parser')
    data = []
    category = None

    for element in soup.find_all(True, {'class': ['tr3', 'tr1', 'tr2']}):
        if 'tr3' in element.get('class', []):
            category = element.get_text(strip=True)
        elif category and ('tr1' in element.get('class', []) or 'tr2' in element.get('class', [])):
            cols = element.find_all('td')
            if len(cols) > 2:
                name = cols[2].get_text(strip=True)
                if name and "Nazwisko Imie" not in name and "WARUNKOWO - BADANIA" not in name:
                    club = cols[6].get_text(strip=True) if len(cols) > 6 else "Unknown Club"
                    data.append((category, name, club))

    return data

def insert_data_to_db(data):
    connection = pymysql.connect(host=db_host, user=db_user, password=db_password, db=db_name, charset='utf8mb4')
    try:
        with connection.cursor() as cursor:
            # Fetch the last race number from the database
            cursor.execute(f"SELECT MAX(numerwyscigu) FROM {table_name}")
            last_race_number = cursor.fetchone()[0] or 0

            # Initialize current race and lane numbers
            current_race_number = last_race_number + 1
            current_lane_number = 1

            # Track the last race number for each competitor
            competitor_last_race = {}

            # Prepare SQL query with race and lane numbers
            sql = f"INSERT INTO {table_name} (kategoria, imie, klub, numerwyscigu, tor) VALUES (%s, %s, %s, %s, %s)"

            # Prepare data with race and lane numbers
            data_with_race_lane = []
            for record in data:
                competitor = record[1]  # Assuming the competitor's name is in the second column

                # Check if the competitor has raced in the last 10 races
                if competitor in competitor_last_race and current_race_number - competitor_last_race[competitor] < 10:
                    continue  # Skip this competitor for now

                # Update the last race number for this competitor
                competitor_last_race[competitor] = current_race_number

                # Append race and lane numbers to the record
                record_with_race_lane = record + (current_race_number, current_lane_number)
                data_with_race_lane.append(record_with_race_lane)

                # Update lane and race numbers for the next competitor
                current_lane_number += 1
                if current_lane_number > 8:  # Assuming there are 8 lanes
                    current_lane_number = 1
                    current_race_number += 1

            # Insert data into the database
            cursor.executemany(sql, data_with_race_lane)
        connection.commit()
    finally:
        connection.close()

# Main execution
url = input("Enter the URL to fetch data from: ")
html_data = fetch_data(url)
if html_data:
    parsed_data = parse_data(html_data)
    if parsed_data:
        insert_data_to_db(parsed_data)
        print("Data inserted successfully.")
    else:
        print("No valid data found to insert.")
