/**
 * Moved from JSON files to TS files, as JSONs are not working in composite projects.
 * @see {@link https://github.com/TypeStrong/ts-loader/issues/905}
 */
export default {
  proto:
    'syntax = "proto3";\n\noption java_multiple_files = true;\noption java_outer_classname = "PokeshopProto";\noption objc_class_prefix = "PKS";\n\npackage pokeshop;\n\nservice Pokeshop {\n  rpc getPokemonList (GetPokemonRequest) returns (GetPokemonListResponse) {}\n  rpc createPokemon (Pokemon) returns (Pokemon) {}\n  rpc importPokemon (ImportPokemonRequest) returns (ImportPokemonRequest) {}\n}\n\nmessage ImportPokemonRequest {\n  int32 id = 1;\n  optional bool isFixed = 2;\n}\n\nmessage GetPokemonRequest {\n  optional int32 skip = 1;\n  optional int32 take = 2;\n  optional bool isFixed = 3;\n}\n\nmessage GetPokemonListResponse {\n  repeated Pokemon items = 1;\n  int32 totalCount = 2;\n}\n\nmessage Pokemon {\n  optional int32 id = 1;\n  string name = 2;\n  string type = 3;\n  bool isFeatured = 4;\n  optional string imageUrl = 5;\n}',
};
