// Copyright 2022 Cisco Systems, Inc. All rights reserved.

/*
Package signal provides structs and functions pertaining to signal handling.
Specifically, for each signal, a Router keeps track of the registered Handler
and the status (i.e. whether the signal is to be ignored or handled). If a
signal is being handled, any subsequent received signals will be handled by the
registered Handler. Else, the signal will be ignored. The signal router
additionally allows the user to adjust the signal handling state and registered
handler. The StartFunc() function starts the signal router and the StopFunc()
function stops the signal route and cleans up any resources.
*/
package signal
