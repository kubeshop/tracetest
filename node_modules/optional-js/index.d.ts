export default class Optional<T> {
    private readonly _value;

    private constructor(value?: T | null | undefined);

    /**
     * Returns an Optional describing the given value, if non-null, otherwise returns an empty Optional.
     *
     * @typeparam T the type of the value
     * @param value the possibly-null value to describe
     * @return an Optional with a present value if the specified value is non-null, otherwise an empty Optional
     */
    static ofNullable<T>(value: T | undefined | null): Optional<T>;

    /**
     * Returns an Optional describing the given non-null value.
     *
     * @typeparam T the type of the value
     * @param value the value to describe, which must be non-null
     * @return an Optional with the value present
     * @throws Error if value is null
     */
    static of<T>(value: T): Optional<T>;

    /**
     * Returns an empty Optional instance. No value is present for this Optional.
     *
     * @return an empty Optional
     */
    static empty<T>(): Optional<T>;

    /**
     * If a value is present in this Optional, returns the value, otherwise throws an Error.
     *
     * @return the non-null value held by this Optional
     * @throws Error if the value is null;
     */
    get(): T;

    /**
     * Return true if there is a value present, otherwise false.
     *
     * @return true if there is a value present, otherwise false
     */
    isPresent(): boolean;

    /**
     * If a value is present, invoke the specified consumer with the value, otherwise do nothing.
     *
     * @param consumer function to be executed if a value is present
     */
    ifPresent(consumer: (value: T) => void): void;

    /**
     * If a value is present, and the value matches the given predicate, return an Optional describing the value,
     * otherwise return an empty Optional.
     *
     * @param predicate A predicate to apply to the value, if present
     * @return an Optional describing the value of this Optional if a value is present and the value matches the given
     * predicate, otherwise an empty Optional
     * @throws Error if the predicate is null
     */
    filter(predicate: (value: T) => boolean): Optional<T>;

    /**
     * If a value is present, apply the provided mapping function to it, and if the result is non-null,
     * return an Optional describing the result. Otherwise return an empty Optional.
     *
     * @typeparam U The type of the result of the mapping function
     * @param mapper a mapping function to apply to the value, if present.
     * @return an Optional describing the result of applying a mapping function to the value of this Optional,
     * if a value is present, otherwise an empty Optional
     * @throws Error if the mapping function is null
     */
    map<U>(mapper: (value: T) => U | undefined | null): Optional<U>;

    /**
     * If a value is present, apply the provided Optional-bearing mapping function to it, return that result,
     * otherwise return an empty Optional. This method is similar to map(Function), but the provided mapper is one whose
     * result is already an Optional, and if invoked, flatMap does not wrap it with an additional Optional.
     *
     * @typeparam U The type parameter to the Optional returned by the mapping function
     * @param mapper a mapping function to apply to the value, if present the mapping function
     * @return the result of applying an Optional-bearing mapping function to the value of this Optional,
     * if a value is present, otherwise an empty Optional
     * @throws Error if the mapping function is null or returns a null result
     */
    flatMap<U>(mapper: (value: T) => Optional<U> | undefined | null): Optional<U>;

    /**
     * If a value is present, returns the value, otherwise returns other.
     *
     * @param other the value to be returned, if no value is present. May be null.
     * @return the value, if present, otherwise other
     */
    orElse(other: T): T;

    /**
     * If a value is present, returns the value, otherwise returns the result produced by the supplying function.
     *
     * @param supplier the supplying function that produces a value to be returned
     * @return the value, if present, otherwise the result produced by the supplying function
     * @throws Error if no value is present and the supplying function is null
     */
    orElseGet(supplier: () => T): T;

    /**
     * If a value is present, returns the value, otherwise throws an exception produced by the exception supplying function.
     *
     * @param exceptionSupplier the supplying function that produces an exception to be thrown
     * @return the value, if present
     * @throws Error if no value is present and the exception supplying function is null
     */
    orElseThrow(exceptionSupplier: () => Error): T;

    /**
     * If a value is present, performs the given action with the value, otherwise performs the given empty-based action.
     *
     * @param action the action to be performed, if a value is present
     * @param emptyAction the empty-based action to be performed, if no value is present
     * @throws if a value is present and the given action is null, or no value is present and the given empty-based action is null.
     */
    ifPresentOrElse(action: (value: T) => void, emptyAction: () => void): void;

    /**
     * If a value is present, returns an Optional describing the value, otherwise returns an Optional produced by the supplying function.
     *
     * @param optionalSupplier the supplying function that produces an Optional to be returned
     * @return returns an Optional describing the value of this Optional, if a value is present,
     * otherwise an Optional produced by the supplying function.
     * @throws Error if the supplying function is null or produces a null result
     */
    or(optionalSupplier: () => Optional<T>): Optional<T>;
}
