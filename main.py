from flask import Flask, request, jsonify
from data import db_session
from flask_cors import CORS
from data.models.user import User
from werkzeug.security import generate_password_hash

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
    )

    db_sess.add(user)
    db_sess.commit()

    return jsonify({"status": "success", "register": data}), 200


if __name__ == "__main__":
    app.run(debug=True)