# model.py
from wtforms import SubmitField, BooleanField, StringField, PasswordField, validators
from flask_wtf import FlaskForm
class TradfriForm(FlaskForm):
  ip = StringField('IP', 
                 [validators.DataRequired()])
  key = StringField('Key', 
                 [validators.DataRequired()])
  submit = SubmitField('Submit')