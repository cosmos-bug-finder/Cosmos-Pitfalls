#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{Binary, Coin, CosmosMsg, Deps, DepsMut, Env, MessageInfo, Response, StdResult, Uint128};
use cosmwasm_std::BankMsg;
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{InstantiateMsg, SudoMsg};
use crate::state::CREATOR;

// version info for migration info
const CONTRACT_NAME: &str = "crates.io:infinite-track-beforesend";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

/// Handling contract instantiation
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    CREATOR.save(deps.storage, &info.sender)?;
    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("creator", info.sender))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: (),
) -> Result<Response, ContractError> {
    Ok(Response::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(_deps: Deps, _env: Env, _msg: ()) -> StdResult<Binary> {
    Ok(Binary::default())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn sudo(deps: DepsMut, _env: Env, msg: SudoMsg) -> Result<Response, ContractError> {
    match msg {
        SudoMsg::TrackBeforeSend { from, to, amount } => {
            // Load the creator address that was stored during instantiation
            let creator = CREATOR.load(deps.storage)?;
            
            // Construct the token factory denomination: factory/{creator_address}/bitcoin
            let factory_denom = format!("factory/{}/bitcoin", creator);
            
            let bitcoin_coin = Coin {
                denom: factory_denom.clone(),
                amount: Uint128::new(1),
            };
            
            let bank_msg1 = BankMsg::Send {
                to_address: to.clone(),
                amount: vec![bitcoin_coin.clone()],
            };
            
            let bank_msg2 = BankMsg::Send {
                to_address: to,
                amount: vec![bitcoin_coin],
            };
            
            let cosmos_msg1 = CosmosMsg::Bank(bank_msg1);
            let cosmos_msg2 = CosmosMsg::Bank(bank_msg2);//AnyMessage Url
            
            Ok(Response::new()
                .add_message(cosmos_msg1)
                .add_message(cosmos_msg2))

        }
        SudoMsg::BlockBeforeSend { .. } => {
            // Minimal gas consumption, just return OK
            Ok(Response::new())
        }
    }
}