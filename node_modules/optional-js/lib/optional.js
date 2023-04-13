function Optional(value) {
    this._value = value;
}

Optional.prototype = {
    get: function get() {
        if (isNull(this._value)) {
            throw new Error('optional is empty');
        }
        return this._value;
    },
    isPresent: function isPresent() {
        return !isNull(this._value);
    },
    ifPresent: function ifPresent(consumer) {
        if (!isNull(this._value)) {
            if (!isFunction(consumer)) {
                throw new Error('consumer is not a function');
            }
            consumer(this._value);
        }
    },
    filter: function filter(predicate) {
        if (!isFunction(predicate)) {
            throw new Error('predicate is not a function');
        }
        if (!isNull(this._value) && predicate(this._value)) {
            return new Optional(this._value);
        }
        return new Optional();
    },
    map: function map(mapper) {
        var mappedValue;

        if (!isFunction(mapper)) {
            throw new Error('mapper is not a function');
        }

        if (isNull(this._value)) {
            return new Optional();
        }

        mappedValue = mapper(this._value);

        return isNull(mappedValue) ? new Optional() : new Optional(mappedValue);
    },
    flatMap: function flatMap(mapper) {
        var flatMappedValue;

        if (!isFunction(mapper)) {
            throw new Error('mapper is not a function');
        }

        if (isNull(this._value)) {
            return new Optional();
        }

        flatMappedValue = mapper(this._value);

        if (isNull(flatMappedValue) || isNull(flatMappedValue.get)) {
            throw new Error('mapper does not return an Optional');
        }

        return flatMappedValue;
    },
    peek: function peek(peeker) {
        if (!isFunction(peeker)) {
            throw new Error('peeker is not a function');
        }

        if (isNull(this._value)) {
            return new Optional();
        }

        peeker(this._value);

        return new Optional(this._value);
    },
    orElse: function orElse(other) {
        return isNull(this._value) ? other : this._value;
    },
    orElseGet: function orElseGet(supplier) {
        if (!isFunction(supplier)) {
            throw new Error('supplier is not a function');
        }
        if (isNull(this._value)) {
            return supplier();
        } else {
            return this._value;
        }
    },
    orElseThrow: function orElseThrow(exceptionSupplier) {
        if (isNull(this._value)) {
            if (!isFunction(exceptionSupplier)) {
                throw new Error('exception provider is not a function');
            }

            throw exceptionSupplier();
        }
        return this._value;
    },
    ifPresentOrElse: function ifPresentOrElse(action, emptyAction) {
        if (!isNull(this._value)) {
            if (!isFunction(action)) {
                throw new Error('action is not a function')
            }
            action(this._value)
        } else {
            if (!isFunction(emptyAction)) {
                throw new Error('emptyAction is not a function')
            }
            emptyAction();
        }
    },
    or: function or(optionalSupplier) {
        if (isNull(this._value)) {
            if (!isFunction(optionalSupplier)) {
                throw new Error('optionalSupplier is not a function')
            }
            return optionalSupplier();
        }
        return this;
    },
    hashCode: function hashMap() {
        // Here just to complete the Java Optional API.
        return -1;
    }
};

function isNull(value) {
    return (value === undefined || value === null);
}

function isFunction(value) {
    return typeof value === 'function';
}

module.exports = Optional;
