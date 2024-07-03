package finance

import (

  //"fmt"
	"github.com/gotk3/gotk3/gtk"
  "gogtk/finance/crypto"
  "gogtk/finance/equities"
  "gogtk/finance/banking"
)

func FinancePage() *gtk.Box{
  
  //fmt.println("\n-----------------------------\n Creating Finance Page \n-----------------------------\n")

  financeBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil  }

  financeStack, err := gtk.StackNew()
  if err != nil {
    return nil
  }

  //--------------------------------------------------

  stackSwitcher, err := gtk.StackSwitcherNew()
  if err != nil {
    return nil
  }

  stackSwitcher.SetStack(financeStack)

  //--------------------------------------------------


  //--------------------------------------------------

  cryptoPage, err := crypto.CryptoPage()
  if err != nil {
    return nil
  }
  financeStack.AddTitled(cryptoPage, "Crypto", "Crypto")

  // ----------------------------------------------------------------------------------------------------
  
  equitiesPage, err := equities.EquitiesPage()
  if err != nil {
    return nil
  }
  financeStack.AddTitled(equitiesPage, "Equities", "Equities")

  // ----------------------------------------------------------------------------------------------------
  
  bankingPage, err := banking.BankingPage()
  if err != nil {
    return nil
  }
  financeStack.AddTitled(bankingPage, "Banking", "Banking")
  


  // ----------------------------------------------------------------------------------------------------
  //
  //
  //

  financeBox.PackStart(stackSwitcher, false, false, 0)
  financeBox.PackStart(financeStack, true, true, 0)

  return financeBox
}

func createSavingsTab() (*gtk.Box, error) {
  box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("SavingsTab")
  if err != nil {
    return nil, err
  }
  box.PackStart(label, false, false, 0)

  return box, nil
}

func createCheckingTab() (*gtk.Box, error) {
  box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("CheckingTab")
  if err != nil {
    return nil, err
  }
  box.PackStart(label, false, false, 0)

  return box, nil
}

func createInvestmentsTab() (*gtk.Box, error) {
  box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("InvestmentsTab")
  if err != nil {
    return nil, err
  }
  box.PackStart(label, false, false, 0)

  return box, nil
}


