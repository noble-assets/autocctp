package autocctp_test

// 1. Packet not transfer type
// 2. Sender chain is source chain
// 3. types.Memo is empty
// 4. Packet is deposit for burn - amount is nil
// 5. Packet is deposit for burn - fee recipient is nil
// 6. Packet is deposit for burn - fee recipient is invalid bech32
// 7. Packet is deposit for burn - amount is invalid
// 8. Packet is deposit for burn - specified amount is greater than packet amount
// 9. Packet is deposit for burn - fee transfer failed
// 10. Packet is deposit for burn - deposit success - assert ack
// 11. Packet is deposit for burn with caller - amount is nil
// 12. Packet is deposit for burn with caller - fee recipient is nil
// 13. Packet is deposit for burn with caller - fee recipient is invalid bech32
// 14. Packet is deposit for burn with caller - amount is invalid
// 15. Packet is deposit for burn with caller - specified amount is greater than packet amount
// 16. Packet is deposit for burn with caller - fee transfer failed
// 17. Packet is deposit for burn with caller - deposit success - assert ack
