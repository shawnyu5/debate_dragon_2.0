from flask import Flask, request

import train

app = Flask("Ivan detector")

@app.route('/', methods = ['POST'])
def hello_world():
    print(f"hello_world request.method: {str(request.method)}") # __AUTO_GENERATED_PRINT_VAR__
    request_data = request.get_data()
    isIvan = train.predict(request_data)[0]
    print(f"hello_world isIvan: {str(isIvan)}") # __AUTO_GENERATED_PRINT_VAR__
    return str(isIvan)
