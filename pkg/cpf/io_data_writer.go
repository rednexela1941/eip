package cpf

func (self *_Writer) AddIODataItem() IWriter {
	w := self.AddItem()
	w.SetTypeID(ConnectedTransportPacket)
	return w
}
