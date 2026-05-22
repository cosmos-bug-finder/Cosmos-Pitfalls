use cosmwasm_std::Addr;
use cw_storage_plus::Item;

pub const CREATOR: Item<Addr> = Item::new("creator");