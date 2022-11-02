export {};

declare global {
  interface Window {
    // TODO include node types
    // IpcTypes
    nodeIpc: any;
  }
}
