# app.py
#!bin/python
from flask import Flask, request, render_template
from forms import TradfriForm
from flask_bs4 import Bootstrap
app = Flask(__name__)

app.config.from_mapping(
    SECRET_KEY=b'\xd6\x04\xbdj\xfe\xed$c\x1e@\xad\x0f\x13,@G')
Bootstrap(app)
app.config['BOOTSTRAP_SERVE_LOCAL'] = True
app.config["BOOTSTRAP_USE_TOASTS"] = False

@app.route('/', methods=['GET', 'POST'])
def registration():
    form = TradfriForm(request.form)
    if request.method == 'POST' and form.validate_on_submit():
        return 'We confirm your registration!'
    return render_template('preferences.html', form=form)

@app.route('/mqtt', methods=['GET', 'POST'])
def page_mqtt():
    form = TradfriForm(request.form)
    if request.method == 'POST' and form.validate_on_submit():
        return 'We confirm your registration!'
    return render_template('mqtt.html', form=form)


if __name__ == '__main__':
    app.run(debug=True)