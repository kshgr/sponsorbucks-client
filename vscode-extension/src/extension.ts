import * as vscode from "vscode";

let statusItem: vscode.StatusBarItem | undefined;
let interval: NodeJS.Timeout | undefined;
let index = 0;

const demoLines = [
  "Sponsored · Deploy APIs faster ↗",
  "Sponsored · Ship logs without noise ↗",
  "Sponsored · Postgres that scales ↗"
];

export function activate(context: vscode.ExtensionContext) {
  statusItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
  statusItem.name = "SponsorBucks";
  statusItem.tooltip = "SponsorBucks sponsored wait-state line. Demo mode.";
  statusItem.command = "sponsorbucks.startDemo";

  context.subscriptions.push(statusItem);

  context.subscriptions.push(
    vscode.commands.registerCommand("sponsorbucks.startDemo", () => {
      startDemo();
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand("sponsorbucks.stopDemo", () => {
      stopDemo();
    })
  );

  statusItem.hide();
}

function startDemo() {
  if (!statusItem) return;
  statusItem.show();
  statusItem.text = "$(megaphone) " + demoLines[index % demoLines.length];

  if (interval) clearInterval(interval);
  interval = setInterval(() => {
    index += 1;
    if (statusItem) {
      statusItem.text = "$(megaphone) " + demoLines[index % demoLines.length];
    }
  }, 5000);

  vscode.window.showInformationMessage("SponsorBucks demo placement started.");
}

function stopDemo() {
  if (interval) {
    clearInterval(interval);
    interval = undefined;
  }
  statusItem?.hide();
  vscode.window.showInformationMessage("SponsorBucks demo placement stopped.");
}

export function deactivate() {
  stopDemo();
}
