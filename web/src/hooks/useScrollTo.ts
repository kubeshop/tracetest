const useScrollTo = () => {
  const scrollTo = (id: string) => {
    const element = document.getElementById(id);

    if (element) {
      element.scrollIntoView({
        behavior: 'smooth',
      });
    }
  };

  return scrollTo;
};

export default useScrollTo;
