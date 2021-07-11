# app.py
#!bin/python
from flask import Flask, request, render_template, Markup
from forms import TradfriForm, MQTTForm, InterfaceForm
from flask_bs4 import Bootstrap

app = Flask(__name__)

app.config.from_mapping(SECRET_KEY=b"\xd6\x04\xbdj\xfe\xed$c\x1e@\xad\x0f\x13,@G")
Bootstrap(app)
app.config["BOOTSTRAP_SERVE_LOCAL"] = True
app.config["BOOTSTRAP_USE_TOASTS"] = False

def render_main_page(active_page="tradfri", status=""):
    form_tradfri = TradfriForm(request.form)
    form_mqtt = MQTTForm(request.form)
    form_interface = InterfaceForm(request.form)

    return render_template(
        "preferences.html",
        tradfri_active="show active" if active_page == "tradfri" else "",
        mqtt_active="show active" if active_page == "mqtt" else "",
        interface_active="show active" if active_page == "interface" else "",
        form_tradfri=Markup(render_template("tradfri.html", form=form_tradfri)),
        form_mqtt=Markup(render_template("mqtt.html", form=form_mqtt)),
        form_interface=Markup(render_template("interface.html", form=form_interface)),
        status = status
    )

@app.route("/", methods=["GET"])
def index():
    return render_main_page()

@app.route("/tradfri", methods=["GET", "POST"])
def tradfri():
    form = TradfriForm(request.form)

    if request.method == 'POST' and form.validate_on_submit():
        return render_main_page(active_page="tradfri", status="Done")
    
    return render_main_page(active_page="tradfri")

@app.route("/mqtt", methods=["GET", "POST"])
def mqtt():
    form = MQTTForm(request.form)

    if request.method == 'POST' and form.validate_on_submit():
        return render_main_page(active_page="mqtt", status="Done")
    
    return render_main_page(active_page="mqtt")

if __name__ == "__main__":
    app.run(debug=True)
