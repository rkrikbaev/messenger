from flask import Flask, render_template, request
import sqlite3
from datetime import datetime

# from flask_sqlalchemy import SQLAlchemy

# Имя базы данных SQLite
database_name = "mydatabase.db"

app = Flask(__name__)

def get_db_connection():
    conn = sqlite3.connect(database_name)
    conn.row_factory = sqlite3.Row
    return conn

@app.route('/health')
def health():
    # Get the current date and time
    current_datetime = datetime.now()
    # Pass the current date and time to the HTML template
    return render_template('index.html', current_datetime=current_datetime)
    
@app.route('/')
def ajax_table():
    return render_template('prod_table.html', title='Product Table')

@app.route('/api/data/<table>')
def data(table):

    start_date = request.args.get('start_date').replace('-','.')
    end_date = request.args.get('end_date').replace('-','.')

    print(start_date)
    print(end_date)

    conn = get_db_connection()
    cursor = conn.cursor()

    query = f'SELECT * FROM {table} WHERE date BETWEEN ? AND ?'
    cursor.execute(query, (start_date, end_date))

    # Fetch all the rows
    data = cursor.fetchall()
    
    print(data)
    cursor.close()

    result = [
                {
                    'date': row[0],
                    'dev1_density': row[1],
                    'dev1_volume': row[2],
                    'dev1_temperature': row[3],
                    'dev1_massflowbegin': row[4],
                    'dev1_massflowend': row[5],
                    'dev1_mass': row[6],
                    'dev1_masstotalizer': row[7],
                    'dev2_density': row[8],
                    'dev2_volume': row[9],
                    'dev2_temperature': row[10],
                    'dev2_level': row[11],
                    'dev2_mass': row[12],
                    'dev3_density': row[13],
                    'dev3_volume': row[14],
                    'dev3_temperature': row[15],
                    'dev3_level': row[16],
                    'dev3_mass': row[17]
                } for row in data
    ]

    return {'data': result}

if __name__ == '__main__':

    app.run(host='0.0.0.0', port=5005, debug=True)
