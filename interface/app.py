# app.py
#!bin/python
from flask import Flask, request, render_template, Markup
from forms import TradfriForm, MQTTForm
from flask_bs4 import Bootstrap

app = Flask(__name__)

app.config.from_mapping(SECRET_KEY=b"\xd6\x04\xbdj\xfe\xed$c\x1e@\xad\x0f\x13,@G")
Bootstrap(app)
app.config["BOOTSTRAP_SERVE_LOCAL"] = True
app.config["BOOTSTRAP_USE_TOASTS"] = False


@app.route("/", methods=["GET"])
def registration():
    form_tradfri = TradfriForm(request.form)
    form_mqtt = MQTTForm(request.form)

    if form_tradfri.submit.data:
        print("Tradfri")

    if form_mqtt.submit.data:
        print("mqtt")

    # if request.method == "POST" and form_tradfri.validate_on_submit():
    #    return "We confirm your registration!"

    return render_template(
        "preferences.html",
        tradfri_active="active",
        mqtt_active="",
        form_tradfri=Markup(render_template("tradfri.html", form=form_tradfri)),
        form_mqtt=Markup(render_template("mqtt.html", form=form_mqtt)),
    )


if __name__ == "__main__":
    app.run(debug=True)
