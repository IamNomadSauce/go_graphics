package banking

import (
  "github.com/gotk3/gotk3/gtk"
  "gogtk/finance/banking/savings"
  "gogtk/finance/banking/checking"

)


func BankingPage() (*gtk.Box, error) {
  bankingBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  bankingStack, err := gtk.StackNew()
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("BankingPage")
  if err != nil {
    return nil, err
  }
  
  // ------------------------------------------------------------

  stackSwitcher, err := gtk.StackSwitcherNew()
  if err != nil {
    return nil, err
  }

  stackSwitcher.SetStack(bankingStack)

  // ------------------------------------------------------------

  savingsPage, err := savings.SavingsPage()
  if err != nil {
    return nil, err
  }

  bankingStack.AddTitled(savingsPage, "Savings", "Savings")

  //--------------------------------------------------
  
  checkingPage, err := checking.CheckingPage()
  if err != nil {
    return nil, err
  }
  bankingStack.AddTitled(checkingPage, "Checking", "checking")

  //--------------------------------------------------
 
  bankingBox.PackStart(label, false, false, 0)
  bankingBox.PackStart(stackSwitcher, false, false, 0)
  bankingBox.PackStart(bankingStack, true, true, 0)

  return bankingBox, nil 
  //--------------------------------------------------
}


