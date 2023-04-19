/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package schemas

// Implements Container interface
type container struct {
	name      string
	schema    QName
	minOccurs Occurs
	maxOccurs Occurs
}

func newContainer(name string, schema QName, minOccurs, maxOccurs Occurs) container {
	return container{
		name:      name,
		schema:    schema,
		minOccurs: minOccurs,
		maxOccurs: maxOccurs,
	}
}

func (cont *container) IsSys() bool { return IsSysContainer(cont.name) }

func (cont *container) MaxOccurs() Occurs { return cont.maxOccurs }

func (cont *container) MinOccurs() Occurs { return cont.minOccurs }

func (cont *container) Name() string { return cont.name }

func (cont *container) Schema() QName { return cont.schema }