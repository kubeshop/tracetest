var Optional = require('./lib/optional.js');

module.exports = {
    empty: function empty() {
        return new Optional();
    },
    of: function of(value) {
        if (value === undefined || value === null) {
            throw new Error('value is not defined');
        }
        return new Optional(value);
    },
    ofNullable: function ofNullable(value) {
        return new Optional(value);
    }
};
