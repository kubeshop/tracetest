const visiblePortion = 100;

export function visiblePortionFuction() {
  return {visiblePortion, height: `calc(100% - ${visiblePortion}px)`};
}
