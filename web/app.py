from flask import Flask, render_template, request, jsonify
import sqlite3 as sql
from datetime import datetime
import os
from flask import send_from_directory


# Имя базы данных SQLite
database_name = "/Users/rustamkrikbayev/projects/parser/web/mydatabase.db"

app = Flask(__name__)

def get_db_connection():
    conn = sql.connect(database_name)
    conn.row_factory = sql.Row
    return conn

@app.route('/health')
def health():
    # Get the current date and time
    current_datetime = datetime.now()
    # Pass the current date and time to the HTML template
    return render_template('index.html', current_datetime=current_datetime)


app = Flask(__name__)

@app.route('/update_data', methods=['POST'])
def update_data():
    if request.method == 'POST':
        try:
            # Retrieve the updated data from the POST request as JSON
            updated_data = request.get_json()
            print('get POST request', updated_data)
            # Extract the relevant data from the JSON request

            record_id = updated_data.get('record_id')
            updated_status = updated_data.get('status')

            print(record_id)
            print(updated_status)

            # Update the "messages" table in the database
            con = sql.connect(database_name)
            cur = con.cursor()

            # Use an SQL UPDATE statement to update the "status" and "sent" columns
            cur.execute("UPDATE messages SET status = ? WHERE id = ?", (updated_status, record_id))
            con.commit()
            con.close()

            response_data = {'messages': 'Data updated successfully'}
            return jsonify(response_data), 200  # Respond with JSON and HTTP status code 200 (OK)

        except Exception as e:
            error_message = {'error': str(e)}
            print(e)
            return jsonify(error_message), 500  # Respond with an error and HTTP status code 500 (Internal Server Error)

    # Handle invalid requests or other conditions here
    return jsonify({'error': 'Invalid request'}), 400  # Respond with an error and HTTP status code 400 (Bad Request)

@app.route('/get_data')
def get_data():
    con = get_db_connection()
    cursor = con.cursor()
    try:
        query = f'SELECT * FROM messages ORDER BY id ASC LIMIT 10'
        # cursor.execute(query, (start_date, end_date))
        cursor.execute(query)

        # Fetch all the rows
        data = cursor.fetchall()
        
        print(data)
        cursor.close()

        data = [
                    {
                        'id': row[0],
                        'date': row[1],
                        'body': row[2],
                        'status': row[3],
                        'created': row[4],
                        'sent': row[5]
                    } for row in data
        ]
    except sql.Error as e:
        print("SQLite Error:", e)
    finally:
        con.close()        
    return  data

@app.route('/favicon.ico')
def favicon():
    return send_from_directory(os.path.join(app.root_path, 'static'),
                               'favicon.ico', mimetype='image/favicon.png')

@app.route('/')
def data():

    # start_date = request.args.get('start_date').replace('-','.')
    # end_date = request.args.get('end_date').replace('-','.')

    # print(start_date)
    # print(end_date)

    con = get_db_connection()
    cursor = con.cursor()
    try:
        query = f'SELECT * FROM messages ORDER BY id ASC LIMIT 10'
        # cursor.execute(query, (start_date, end_date))
        cursor.execute(query)

        # Fetch all the rows
        data = cursor.fetchall()
        
        print(data)
        cursor.close()

        data = [
                    {
                        'id': row[0],
                        'date': row[1],
                        'body': row[2],
                        'status': row[3],
                        'created': row[4],
                        'sent': row[5]
                    } for row in data
        ]

        return render_template('production.html', title='Журнал регистрации', datas=data)
    except sql.Error as e:
        print("SQLite Error:", e)
    finally:
        con.close()

if __name__ == '__main__':

    app.run(host='0.0.0.0', port=5005, debug=True)
