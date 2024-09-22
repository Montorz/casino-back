from flask import Flask, request, jsonify
from data import db_session
from flask_cors import CORS
from data.models.user import User

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "http://127.0.0.1:5500"}})  # Разрешаем только определённый источник
db_session.global_init('db/casino.db')


@app.route('/register', methods=['POST', 'OPTIONS', 'GET'])
def register():
    if request.method == 'OPTIONS':
        return jsonify({'status': 'ok'}), 200  # Ответ на preflight запрос

    data = request.get_json()
    user = User(
        name=data['name'],
        password=data['password'],
    )

    db_sess = db_session.create_session()
    db_sess.add(user)
    db_sess.commit()

    return jsonify({"status": "success", "register": data}), 200


if __name__ == "__main__":
    app.run(debug=True)