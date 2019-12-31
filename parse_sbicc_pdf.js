'use strict';

const PdfReader = require('pdfreader').PdfReader;
const moment = require('moment');
const parseArgs = require('minimist');


let fnParseData = (arrData) => {
    let finalCleanData = [];

    arrData.forEach(element => {
        if(element.length > 3) {
            let indexDate = 0;
            let indexAmount = element.length - 1;
            let indexType = element.length - 2;
            // make sure 1st element is a date
            let transDate = moment(element[indexDate], "DD/MM/YYYY");
            if(transDate.isValid()) {
                transDate = transDate.format("DD/MM/YYYY");

                // amount should be greater then 0
                let transAmount = element[indexAmount];
                transAmount = parseFloat(transAmount, 10)
                if(transAmount > 0) {
                    let transType = element[indexType];
                    transType = transType.toLowerCase().trim();
                    if(transType == 'debit' || transType == 'credit') {
                        element.splice(0, 1)
                        element.splice(element.length - 1, 1)
                        element.splice(element.length - 1, 1)
                        finalCleanData.push({
                            date: transDate,
                            type: transType,
                            amount: transAmount,
                            description: element.join(' ')
                        });
                    } else {
                        //console.log('Bad data 1: ' + element)
                    }

                } else {
                    //console.log('Bad data 2: ' + element)
                }

            } else {
                //console.log('Bad data 3: ' + element)
            }
        }
    });

    console.log(JSON.stringify(finalCleanData))
}

let fnParsePDF = (filePath) => {
    let prevY = '';
    let thisText = [];
    let thisData = [];

    new PdfReader().parseFileItems(filePath, function(err, item) {
        if (err) {
            throw err;
        }
        else if (!item) {
            fnParseData(thisData);
        }
        else if (item.y) {
            let diff = item.y - prevY;
            if((item.y == prevY) || (diff < 0.4 && diff > -4)) {
                thisText.push(item.text);
            } else {
                thisText = [item.text];
                thisData.push(thisText);
                prevY = item.y
            }
        }
    });
}

const argv = parseArgs(process.argv.slice(2));

let filePath = argv['_'][0];
fnParsePDF(filePath)