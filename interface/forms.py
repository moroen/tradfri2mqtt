# model.py
from wtforms import SubmitField, BooleanField, StringField, PasswordField, validators, BooleanField
from flask_wtf import FlaskForm


class TradfriForm(FlaskForm):
    ip = StringField("IP", [validators.DataRequired()])
    key = StringField("Key", [validators.DataRequired()])
    submit = SubmitField("Submit")


class MQTTForm(FlaskForm):
    host = StringField("Broker IP", [validators.DataRequired()])
    port = StringField("Broker Port", [validators.DataRequired()], default="1883")
    submit = SubmitField("Submit")

class InterfaceForm(FlaskForm):
    enabled = BooleanField("Enable http interface")
    submit = SubmitField("Submit")