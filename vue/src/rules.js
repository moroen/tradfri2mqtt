export const stringRules = [val => (val && val.length > 0) || 'Field can not be empty'] 

export const timerRule = [v => (!isNaN(parseFloat(v)) && v >= 0 && v <= 3600) || 'Please enter an integer between 0 and 3600']

export const portRule = [v => (!isNaN(parseFloat(v)) && v >= 0 && v <= 65535) || 'Please enter an integer between 0 and 65535']
