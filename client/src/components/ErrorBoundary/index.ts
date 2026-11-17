import { Component, type PropsWithChildren, type ReactNode } from 'react';

export type ErrorBoundaryProps = PropsWithChildren<{
  fallback: (error: unknown) => ReactNode;
  errorKey?: unknown;
}>;

type ErrorBoundaryState = {
  hasError: boolean;
  error: unknown;
  errorKey: unknown;
};

export default class ErrorBoundary extends Component<
  ErrorBoundaryProps,
  ErrorBoundaryState
> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: undefined,
      errorKey: props.errorKey,
    };
  }

  static getDerivedStateFromProps(
    props: ErrorBoundaryProps,
    state: ErrorBoundaryState,
  ) {
    if (state.errorKey === props.errorKey) return null;
    return {
      hasError: false,
      error: undefined,
      errorKey: props.errorKey,
    };
  }

  static getDerivedStateFromError(error: unknown) {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback(this.state.error);
    }

    return this.props.children;
  }
}
