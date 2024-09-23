import datetime
from flask import Flask, request, jsonify
from data import db_session
from flask_cors import CORS
from data.models.user import User
from werkzeug.security import generate_password_hash, check_password_hash

app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "http://127.0.0.1:5500"}})
db_session.global_init('db/casino.db')


@app.route('/register', methods=['POST', 'OPTIONS', 'GET'])
def register():
    if request.method == 'OPTIONS':
        return jsonify({'status': 'ok'}), 200

    data = request.get_json()
    db_sess = db_session.create_session()

    # Проверка на пользователя с таким же логином
    if db_sess.query(User).filter(User.login == data['login']).first():
        return jsonify({"status": "error", "message": "Пользователь с таким логином уже существует"}), 400

    user = User(
        name=data['name'],
        login=data['login'],
        hashed_password=generate_password_hash(data['password']),
        created_date=datetime.datetime.now().strftime('%d-%m-%Y'),
    )

    db_sess.add(user)
    db_sess.commit()

    return jsonify({"status": "success", "register": data}), 200

@app.route('/login', methods=['POST', 'OPTIONS', 'GET'])
def login():
    if request.method == 'OPTIONS':
        return jsonify({'status': 'ok'}), 200

    data = request.get_json()
    db_sess = db_session.create_session()

    # Проверка пользователя на правильность ввода логина и пароля
    user = db_sess.query(User).filter(User.login == data['login']).first()

    if not user or not check_password_hash(user.hashed_password, data['password']):
        return jsonify({"status": "error", "message": "Неправильный логин или пароль"}), 400

    return jsonify({"status": "success", "login": data}), 200

@app.route('/account', methods=['POST', 'OPTIONS', 'GET'])
def account():
    if request.method == 'OPTIONS':
        return jsonify({'status': 'ok'}), 200

    data = request.get_json()
    db_sess = db_session.create_session()

    user = db_sess.query(User).filter(User.login == data['login']).first()

    if not user:
        return jsonify({"status": "error", "message": "Пользователь не найден"}), 404

    # Формирование данных пользователя для ответа
    user_data = {
        "id": user.id,
        "name": user.name,
        "login": user.login,
        "created_date": user.created_date,
    }

    return jsonify({"status": "success", "account": user_data}), 200

if __name__ == "__main__":
    app.run(debug=True)